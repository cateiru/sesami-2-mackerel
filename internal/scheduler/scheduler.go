package scheduler

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cateiru/sesami-2-mackerel/internal/mackerel"
	"github.com/cateiru/sesami-2-mackerel/internal/sesami"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron           *cron.Cron
	sesamiClient   *sesami.Client
	mackerelClient *mackerel.Client
}

func New() *Scheduler {
	return &Scheduler{
		cron:           cron.New(),
		sesamiClient:   sesami.NewClient(),
		mackerelClient: mackerel.NewClient(),
	}
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
}

func (s *Scheduler) dailyTask() {
	log.Printf("定期タスクを実行中... [%s]", time.Now().Format("2006-01-02 15:04:05"))

	status, err := s.sesamiClient.GetDeviceStatus()
	if err != nil {
		log.Printf("SESAMI API呼び出しエラー: %v", err)
		return
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
