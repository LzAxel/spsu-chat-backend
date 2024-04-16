package main

import (
	"spsu-chat/internal/app"
	"spsu-chat/internal/config"
)

func main() {
	cfg := config.ReadConfig()

	app := app.New(cfg)

	app.Start()
}
