package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	startDiscordBot()

	keepAlive("Keeping Alive")
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
