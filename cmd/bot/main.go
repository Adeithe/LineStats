package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"LineStats/internal/app/command"
	"LineStats/internal/app/discord"
	"LineStats/internal/app/twitch"
	"LineStats/internal/pkg/postgres"
	"LineStats/internal/pkg/prometheus"

	_ "github.com/joho/godotenv/autoload"
)

var sc chan os.Signal

func main() {
	rand.Seed(time.Now().Unix())
	sc = make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"), os.Getenv("POSTGRES_DB"),
	)
	if err := postgres.Connect(dbInfo); err != nil {
		panic(err)
	}
	defer postgres.Close()

	go prometheus.Init()

	command.Init()
	twitch.New()
	twitch.Start(os.Getenv("LINESTATS_TWITCH_USERNAME"), os.Getenv("LINESTATS_TWITCH_TOKEN"))

	discord.New()
	discord.Start(os.Getenv("LINESTATS_DISCORD_CLIENT"), os.Getenv("LINESTATS_DISCORD_TOKEN"))

	go ticker()

	<-sc
}

func ticker() {
	for {
		time.Sleep(30 * time.Second)
		if err := postgres.Ping(); err != nil {
			fmt.Println(err)
			break
		}
	}
	close()
}

func close() {
	if sc == nil {
		panic("attempt to close before ready")
	}
	sc <- syscall.SIGINT
}
