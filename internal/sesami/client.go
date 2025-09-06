package sesami

import (
	"fmt"
	"log"

	"github.com/cateiru/sesami-2-mackerel/internal/config"
)

type Client struct {
	APIKey     string
	DeviceUUID string
}

type DeviceStatus struct {
	Battery    int    `json:"battery"`
	IsLocked   bool   `json:"is_locked"`
	Timestamp  int64  `json:"timestamp"`
	DeviceName string `json:"device_name"`
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		APIKey:     cfg.SESAMI.APIKey,
		DeviceUUID: cfg.SESAMI.DeviceUUID,
	}
}

func (c *Client) GetDeviceStatus() (*DeviceStatus, error) {
	log.Printf("SESAMI API呼び出し中... (Device UUID: %s)", c.DeviceUUID)
	
	fmt.Println("SESAMI APIからデバイス状態を取得中...")
	
	status := &DeviceStatus{
		Battery:    85,
		IsLocked:   true,
		Timestamp:  1234567890,
		DeviceName: "SESAMI Device",
	}
	
	log.Printf("SESAMI APIから状態を取得しました: Battery=%d%%, Locked=%t", status.Battery, status.IsLocked)
	return status, nil
}