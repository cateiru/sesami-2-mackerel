package mackerel

import (
	"log"

	"github.com/cateiru/sesami-2-mackerel/internal/config"
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

func NewClient(cfg *config.Config) *Client {
	return &Client{
		APIKey: cfg.Mackerel.APIKey,
	}
}

func (c *Client) SendMetrics(status *sesami.DeviceStatus) error {
	log.Printf("Mackerelにメトリクスを送信中...")

	return nil
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
