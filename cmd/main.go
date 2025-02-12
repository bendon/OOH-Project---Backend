package main

import (
	"bbscout/config/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	server := server.NewApiServer(":" + getEnv("PORT", "8700"))
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
