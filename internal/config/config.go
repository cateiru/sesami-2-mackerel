package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	SESAMI struct {
		APIKey     string
		DeviceUUID string
	}
	Mackerel struct {
		APIKey string
	}
}

func Load() *Config {
	cfg := &Config{}

	cfg.SESAMI.APIKey = getEnv("SESAMI_API_KEY", true)
	cfg.SESAMI.DeviceUUID = getEnv("SESAMI_DEVICE_UUID", true)
	cfg.Mackerel.APIKey = getEnv("MACKEREL_API_KEY", true)

	log.Printf("設定を読み込みました: SESAMI Device UUID=%s", cfg.SESAMI.DeviceUUID)
	return cfg
}

func getEnv(key string, required bool) string {
	value := os.Getenv(key)

	if required && value == "" {
		log.Fatalf("必須の環境変数 %s が設定されていません", key)
	}

	if value == "" {
		log.Printf("環境変数 %s は設定されていません（オプション）", key)
	}

	return value
}

func (c *Config) Validate() error {
	if c.SESAMI.APIKey == "" {
		return fmt.Errorf("SESAMI_API_KEY が設定されていません")
	}

	if c.SESAMI.DeviceUUID == "" {
		return fmt.Errorf("SESAMI_DEVICE_UUID が設定されていません")
	}

	if c.Mackerel.APIKey == "" {
		return fmt.Errorf("MACKEREL_API_KEY が設定されていません")
	}

	return nil
}
