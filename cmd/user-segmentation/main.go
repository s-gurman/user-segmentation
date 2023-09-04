package main

import (
	"flag"
	"log"

	"github.com/s-gurman/user-segmentation/config"
	"github.com/s-gurman/user-segmentation/internal/app"
)

var configPath = flag.String("config", "./config/config.yml", "App config path.")

func main() {
	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		log.Panicf("main - new config: %s", err)
	}

	app.Run(cfg)
}
