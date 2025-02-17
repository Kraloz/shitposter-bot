package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"shitposter-bot/database"
	"shitposter-bot/discord"
	"shitposter-bot/tenor"
	"shitposter-bot/twitter"

	"github.com/joho/godotenv"
)

func init() {
  // loads values from .env into the system
  if err := godotenv.Load(); err != nil {
    log.Print("No .env file found")
  }
}

func main() {
  discord_token := os.Getenv("DISCORD_TOKEN")
  database_path := os.Getenv("DB_PATH")
  tw_access_token := os.Getenv("TW_ACCESS_TOKEN")
  tw_access_token_secret := os.Getenv("TW_ACCESS_TOKEN_SECRET")
  tw_consumer_key := os.Getenv("TW_CONSUMER_KEY")
  tw_consumer_key_secret := os.Getenv("TW_CONSUMER_KEY_SECRET=")
  tenor_token := os.Getenv("TENOR_TOKEN")

	if database_path == "" {
		fmt.Println("Missing Database Path")
	}

	if discord_token == "" || tw_access_token == "" || tw_access_token_secret == "" || tw_consumer_key == "" || tw_consumer_key_secret == "" || tenor_token == "" {
		fmt.Println("Missing tokens")
		return
	}

	database.Start(database_path)
	tenor.Start(tenor_token)
	go twitter.Start(tw_access_token, tw_access_token_secret, tw_consumer_key, tw_consumer_key_secret)
	go discord.Start(discord_token)

	//wait until we want to stop the program
	chnl := make(chan os.Signal, 1)
	signal.Notify(chnl, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-chnl

	//stop and close
	discord.Stop()
	twitter.Stop()
	database.Close()
}
