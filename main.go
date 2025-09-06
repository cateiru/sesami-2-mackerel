package main

import (
	"fmt"

	"github.com/cateiru/sesami-2-mackerel/internal/scheduler"
)

func main() {
	fmt.Println("SESAMI to Mackerel監視プログラムを開始します...")

	s := scheduler.New()
	s.Start()
}
