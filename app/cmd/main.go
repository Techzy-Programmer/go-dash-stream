package main

import (
	"fmt"
	"os"
	"pi-go/internals/handler"
	"time"

	"github.com/joho/godotenv"
)

var isSleeping = false

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	channelID := os.Getenv("CHANNEL_ID")

	go handler.StartFFMPEGThread()
	go handler.StartStreamServer("6619")
	var tgBot = handler.NewTelegramBot(botToken, channelID)

	handler.SubscribeForDistance(func(dist float64) {
		if dist < 10 && !isSleeping {
			go func() {
				tgBot.SendMessage("<b>Warning</b>: Intruder detected!\n<b>Time</b>: " + time.Now().Format(time.RFC1123))
				isSleeping = true

				time.Sleep(20 * time.Second)
				fmt.Println("Waking up...")
				isSleeping = false
			}()
		}
	})

	select {}
}
