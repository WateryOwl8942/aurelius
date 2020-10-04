package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	discord, err := discordgo.New("Bot " + "authentication token")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(discord)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
