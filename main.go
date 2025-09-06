package main

import (
	"fmt"
	"log"

	"github.com/cateiru/sesami-2-mackerel/internal/config"
	"github.com/cateiru/sesami-2-mackerel/internal/scheduler"
)

func main() {
	fmt.Println("SESAMI to Mackerel監視プログラムを開始します...")

	cfg := config.Load()

	s, err := scheduler.New(cfg)
	if err != nil {
		log.Fatalf("スケジューラー初期化エラー: %v", err)
	}
	s.Start()
}
