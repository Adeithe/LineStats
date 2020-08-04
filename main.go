package main

import (
	"LineStats/command/handlers"
	"LineStats/manager"
	"LineStats/postgres"
	"LineStats/twitch"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	twitchUsername = os.Getenv("LINESTATS_TWITCH_USERNAME")
	twitchToken    = os.Getenv("LINESTATS_TWITCH_OAUTH")

	discordClientID = os.Getenv("LINESTATS_DISCORD_CLIENT")
	discordToken    = os.Getenv("LINESTATS_DISCORD_TOKEN")
)

var sc chan os.Signal

func main() {
	sc = make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"), os.Getenv("POSTGRES_DB"),
	)
	if err := postgres.Connect(psqlInfo); err != nil {
		panic(err)
	}
	defer postgres.Close()
	if err := postgres.CreateTables(); err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 30)
			if err := postgres.Ping(); err != nil {
				fmt.Println(err)
				break
			}
		}
		close()
	}()

	manager := manager.New()
	manager.Twitch.SetLogin(twitchUsername, twitchToken)
	manager.Discord.SetLogin(discordClientID, discordToken)

	if err := manager.Connect(); err != nil {
		panic(err)
	}
	handlers.Init()
	manager.Twitch.Join(twitch.ToChannelName(twitchUsername), true)

	go func() {
		for {
			channels, _ := postgres.FetchChannels()
			for _, channel := range channels {
				manager.Twitch.Join(channel.Name, channel.Status)
			}
			time.Sleep(time.Minute * 5)
		}
	}()

	<-sc
	manager.Close()
}

func close() {
	if sc == nil {
		panic("attempt to close before ready")
	}
	sc <- syscall.SIGINT
}
