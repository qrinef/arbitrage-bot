package main

import "github.com/qrinef/arbitrage-bot/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
