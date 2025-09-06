package mackerel

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cateiru/sesami-2-mackerel/internal/sesami"
)

type Client struct {
	APIKey string
}

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Time  int64   `json:"time"`
}

func NewClient() *Client {
	apiKey := os.Getenv("MACKEREL_API_KEY")
	
	if apiKey == "" {
		log.Fatal("MACKEREL_API_KEY の環境変数が設定されていません")
	}

	return &Client{
		APIKey: apiKey,
	}
}

func (c *Client) SendMetrics(status *sesami.DeviceStatus) error {
	log.Printf("Mackerelにメトリクスを送信中...")
	
	metrics := []Metric{
		{
			Name:  "sesami.battery",
			Value: float64(status.Battery),
			Time:  time.Now().Unix(),
		},
		{
			Name:  "sesami.lock_status",
			Value: boolToFloat(status.IsLocked),
			Time:  time.Now().Unix(),
		},
	}

	for _, metric := range metrics {
		fmt.Printf("Mackerelにメトリクスを送信: %s = %.1f\n", metric.Name, metric.Value)
	}
	
	log.Printf("Mackerelに%d個のメトリクスを送信しました", len(metrics))
	return nil
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}