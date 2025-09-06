package scheduler

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cateiru/sesami-2-mackerel/internal/config"
	"github.com/cateiru/sesami-2-mackerel/internal/database"
	"github.com/cateiru/sesami-2-mackerel/internal/mackerel"
	"github.com/cateiru/sesami-2-mackerel/internal/sesami"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron           *cron.Cron
	sesamiClient   *sesami.Client
	mackerelClient *mackerel.Client
	dbClient       *database.Client
}

func New(cfg *config.Config) (*Scheduler, error) {
	dbClient, err := database.NewClient(cfg.Database.Path)
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		cron:           cron.New(),
		sesamiClient:   sesami.NewClient(cfg),
		mackerelClient: mackerel.NewClient(cfg),
		dbClient:       dbClient,
	}, nil
}

func (s *Scheduler) Start() {
	_, err := s.cron.AddFunc("@every 24h", s.dailyTask)
	if err != nil {
		log.Fatalf("cronジョブの追加に失敗しました: %v", err)
	}

	s.cron.Start()
	log.Println("cronスケジューラーが開始されました。24時間間隔でタスクが実行されます。")

	s.dailyTask()

	s.waitForShutdown()
}

func (s *Scheduler) Stop() {
	log.Println("プログラムを終了します...")
	s.cron.Stop()
	if s.dbClient != nil {
		if err := s.dbClient.Close(); err != nil {
			log.Printf("データベース接続クローズエラー: %v", err)
		}
	}
}

func (s *Scheduler) dailyTask() {
	log.Printf("定期タスクを実行中... [%s]", time.Now().Format("2006-01-02 15:04:05"))

	status, err := s.sesamiClient.GetDeviceStatus()
	if err != nil {
		log.Printf("SESAMI API呼び出しエラー: %v", err)
		return
	}

	err = s.dbClient.InsertDeviceStatus(status)
	if err != nil {
		log.Printf("データベース保存エラー: %v", err)
	}

	history, err := s.sesamiClient.GetDeviceHistory()
	if err != nil {
		log.Printf("SESAMI履歴取得エラー: %v", err)
	} else {
		deviceUUID := s.sesamiClient.DeviceUUID
		err = s.dbClient.InsertDeviceHistory(deviceUUID, history)
		if err != nil {
			log.Printf("履歴データベース保存エラー: %v", err)
		} else {
			log.Printf("履歴データを%d件保存しました", len(history))
		}
	}

	err = s.mackerelClient.SendMetrics(status)
	if err != nil {
		log.Printf("Mackerelメトリクス送信エラー: %v", err)
		return
	}

	log.Println("定期タスクが完了しました")
}

func (s *Scheduler) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.Stop()
}
