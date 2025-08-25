package main

import (
	"github.com/hmluck83/txlens-srv/api"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	router := api.NewRouter()
	router.Run()
}
