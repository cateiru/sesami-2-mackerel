package sesami

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cateiru/sesami-2-mackerel/internal/config"
)

type Client struct {
	APIKey       string
	APIUserAgent string
	APITimeout   time.Duration
	DeviceUUID   string
}

type DeviceStatus struct {
	BatteryPercentage int     `json:"batteryPercentage"`
	BatteryVoltage    float64 `json:"batteryVoltage"`
	Position          int     `json:"position"`
	CHSesame2Status   string  `json:"CHSesame2Status"`
	Timestamp         int64   `json:"timestamp"`
}

type HistoryEntry struct {
	Type      string `json:"type"`
	TimeStamp int64  `json:"timeStamp"`
	Tag       string `json:"tag"`
}

type DeviceHistoryResponse struct {
	History []HistoryEntry `json:"history"`
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		APIKey:       cfg.SESAMI.APIKey,
		APIUserAgent: cfg.SESAMI.APIUserAgent,
		APITimeout:   cfg.SESAMI.APITimeout,
		DeviceUUID:   cfg.SESAMI.DeviceUUID,
	}
}

func (c *Client) GetDeviceStatus() (*DeviceStatus, error) {
	url := fmt.Sprintf("https://app.candyhouse.co/api/sesame2/%s", c.DeviceUUID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("User-Agent", c.APIUserAgent)

	client := &http.Client{
		Timeout: c.APITimeout,
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

	return &status, nil
}

func (c *Client) GetDeviceHistory() ([]HistoryEntry, error) {
	url := fmt.Sprintf("https://app.candyhouse.co/api/sesame2/%s/history", c.DeviceUUID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("履歴取得リクエスト作成エラー: %w", err)
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("User-Agent", c.APIUserAgent)

	client := &http.Client{
		Timeout: c.APITimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("履歴取得HTTP リクエストエラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("履歴取得API エラー: ステータスコード %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("履歴レスポンス読み取りエラー: %w", err)
	}

	var historyResp DeviceHistoryResponse
	if err := json.Unmarshal(body, &historyResp); err != nil {
		return nil, fmt.Errorf("履歴JSON パースエラー: %w", err)
	}

	return historyResp.History, nil
}
