package main

import (
	"github.com/trenchesdeveloper/go-store-app/config"
	"github.com/trenchesdeveloper/go-store-app/internal/api"
	"strconv"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	server := api.NewServer(cfg)
	port, err := strconv.Atoi(cfg.ServerPort)

	if err != nil {
		panic(err)
	}
	server.Start(port)
}
