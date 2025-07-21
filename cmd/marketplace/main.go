package main

import (
	"fmt"
	"marketplace/internal/app"
)

func main() {

	market := app.NewApp()

	if err := market.Run(); err != nil {
		fmt.Printf("failed running app: %s", err)
		return
	}
}
