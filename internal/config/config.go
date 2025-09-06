package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	SESAMI struct {
		APIKey       string
		DeviceUUID   string
		APIUserAgent string
		APITimeout   time.Duration
	}
	Mackerel struct {
		APIKey       string
		APIUserAgent string
		APITimeout   time.Duration
		ServiceName  string
	}
	Database struct {
		Path string
	}
	TempFilePath string
}

func Load() *Config {
	cfg := &Config{}

	cfg.SESAMI.APIKey = getEnv("SESAMI_API_KEY", true)
	cfg.SESAMI.DeviceUUID = getEnv("SESAMI_DEVICE_UUID", true)
	cfg.Mackerel.APIKey = getEnv("MACKEREL_API_KEY", true)

	cfg.SESAMI.APIUserAgent = "sesami-2-mackerel/1.0"
	cfg.Mackerel.APIUserAgent = "sesami-2-mackerel/1.0"

	cfg.SESAMI.APITimeout = 30 * time.Second
	cfg.Mackerel.APITimeout = 30 * time.Second

	cfg.Database.Path = getEnv("DATABASE_PATH", true)

	cfg.Mackerel.ServiceName = "sesami2mackerel"

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

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

	if c.Database.Path == "" {
		return fmt.Errorf("DATABASE_PATH が設定されていません")
	}

	return nil
}
