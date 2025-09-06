package mackerel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cateiru/sesami-2-mackerel/internal/config"
	"github.com/cateiru/sesami-2-mackerel/internal/sesami"
)

type Client struct {
	APIKey     string
	UserAgent  string
	Timeout    time.Duration
	httpClient *http.Client
}

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Time  int64   `json:"time"`
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		APIKey:    cfg.Mackerel.APIKey,
		UserAgent: cfg.Mackerel.APIUserAgent,
		Timeout:   cfg.Mackerel.APITimeout,
		httpClient: &http.Client{
			Timeout: cfg.Mackerel.APITimeout,
		},
	}
}

func (c *Client) SendMetrics(status *sesami.DeviceStatus, serviceName string) error {
	log.Printf("Mackerelにメトリクスを送信中...")

	metrics := []Metric{
		{
			Name:  "sesami-battery-percentage",
			Value: float64(status.BatteryPercentage),
			Time:  status.Timestamp,
		},
	}

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("メトリクスJSON作成エラー: %v", err)
	}

	url := fmt.Sprintf("https://api.mackerelio.com/api/v0/services/%s/tsdb", serviceName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %v", err)
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("API呼び出しエラー: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("メトリクス送信失敗: ステータスコード %d", resp.StatusCode)
	}

	log.Printf("メトリクス送信完了: バッテリー残量=%d%%", status.BatteryPercentage)
	return nil
}

type Service struct {
	Name  string   `json:"name"`
	Memo  string   `json:"memo"`
	Roles []string `json:"roles"`
}

type ServiceListResponse struct {
	Services []Service `json:"services"`
}

type ServiceCreateRequest struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

func (c *Client) ListServices() ([]Service, error) {
	req, err := http.NewRequest("GET", "https://api.mackerelio.com/api/v0/services", nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %v", err)
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API呼び出しエラー: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API呼び出し失敗: ステータスコード %d", resp.StatusCode)
	}

	var serviceList ServiceListResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceList); err != nil {
		return nil, fmt.Errorf("レスポンス解析エラー: %v", err)
	}

	return serviceList.Services, nil
}

func (c *Client) CreateService(name string) (*Service, error) {
	reqBody := ServiceCreateRequest{
		Name: name,
		Memo: "",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("JSON作成エラー: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.mackerelio.com/api/v0/services", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %v", err)
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API呼び出しエラー: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("サービス作成失敗: ステータスコード %d", resp.StatusCode)
	}

	var service Service
	if err := json.NewDecoder(resp.Body).Decode(&service); err != nil {
		return nil, fmt.Errorf("レスポンス解析エラー: %v", err)
	}

	return &service, nil
}

func (c *Client) EnsureService(serviceName string) error {
	services, err := c.ListServices()
	if err != nil {
		return fmt.Errorf("サービス一覧取得エラー: %v", err)
	}

	for _, service := range services {
		if service.Name == serviceName {
			log.Printf("サービス '%s' は既に存在します", serviceName)
			return nil
		}
	}

	log.Printf("サービス '%s' が存在しません。作成します...", serviceName)
	service, err := c.CreateService(serviceName)
	if err != nil {
		return fmt.Errorf("サービス作成エラー: %v", err)
	}

	log.Printf("サービス '%s' を作成しました", service.Name)
	return nil
}
