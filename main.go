package main

import (
	"log"

	"github.com/cateiru/sesami-2-mackerel/internal/config"
	"github.com/cateiru/sesami-2-mackerel/internal/scheduler"
)

func main() {
	cfg := config.Load()

	s, err := scheduler.New(cfg)
	if err != nil {
		log.Fatalf("スケジューラー初期化エラー: %v", err)
	}
	s.Start()
}
