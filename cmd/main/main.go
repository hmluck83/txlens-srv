package main

import (
	"fmt"
	"os"

	"github.com/hmluck83/txlens-srv/api"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	fmt.Println(os.Getwd())
	fmt.Println(os.Getenv("GEMINIAPI"))
	router := api.NewRouter()
	router.Run()
}
