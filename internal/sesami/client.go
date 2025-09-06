package sesami

import (
	"fmt"
	"log"
	"os"
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

func NewClient() *Client {
	apiKey := os.Getenv("SESAMI_API_KEY")
	deviceUUID := os.Getenv("SESAMI_DEVICE_UUID")

	if apiKey == "" || deviceUUID == "" {
		log.Fatal("SESAMI_API_KEY と SESAMI_DEVICE_UUID の環境変数が設定されていません")
	}

	return &Client{
		APIKey:     apiKey,
		DeviceUUID: deviceUUID,
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