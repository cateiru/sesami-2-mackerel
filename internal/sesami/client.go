package sesami

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cateiru/sesami-2-mackerel/internal/config"
)

type Client struct {
	APIKey     string
	DeviceUUID string
}

type DeviceStatus struct {
	BatteryPercentage int     `json:"batteryPercentage"`
	BatteryVoltage    float64 `json:"batteryVoltage"`
	Position          int     `json:"position"`
	CHSesame2Status   string  `json:"CHSesame2Status"`
	Timestamp         int64   `json:"timestamp"`
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		APIKey:     cfg.SESAMI.APIKey,
		DeviceUUID: cfg.SESAMI.DeviceUUID,
	}
}

func (c *Client) GetDeviceStatus() (*DeviceStatus, error) {
	url := fmt.Sprintf("https://app.candyhouse.co/api/sesame2/%s", c.DeviceUUID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("X-Api-Key", c.APIKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP リクエストエラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API エラー: ステータスコード %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み取りエラー: %w", err)
	}

	var status DeviceStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("JSON パースエラー: %w", err)
	}

	log.Printf("SESAMI APIから状態を取得しました: Battery=%d%%, Status=%s",
		status.BatteryPercentage, status.CHSesame2Status)

	return &status, nil
}

func (s *DeviceStatus) GetBattery() int {
	return s.BatteryPercentage
}

func (s *DeviceStatus) IsLocked() bool {
	return s.CHSesame2Status == "locked"
}
