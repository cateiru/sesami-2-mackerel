package main

import (
	"fmt"

	"github.com/cateiru/sesami-2-mackerel/internal/config"
	"github.com/cateiru/sesami-2-mackerel/internal/scheduler"
)

func main() {
	fmt.Println("SESAMI to Mackerel監視プログラムを開始します...")

	cfg := config.Load()

	s := scheduler.New(cfg)
	s.Start()
}
