package main

import (
	"github.com/trenchesdeveloper/go-store-app/config"
	"github.com/trenchesdeveloper/go-store-app/internal/api"
)


func main(){
	cfg, err := config.SetupEnv()
	if err != nil {
		panic(err)
	}
	api.StartServer(cfg)
}