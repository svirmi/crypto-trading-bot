package main

import (
	"fmt"

	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
)

func main() {
	fmt.Println("This is main.go!")
	binance.Bhttp()
	binance.Bws()
}
