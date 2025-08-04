package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mbtamuli/minecraft/discord-bot/internal/config"
	"github.com/mbtamuli/minecraft/discord-bot/internal/discord"
	"github.com/mbtamuli/minecraft/discord-bot/internal/docker"
)

func main() {
	cfg := config.Load()

	dockerManager := docker.NewManager(cfg.ComposePath)
	defer dockerManager.Close()

	bot, err := discord.NewBot(cfg.Token, cfg.AllowedRoleID, dockerManager)
	if err != nil {
		panic(fmt.Sprintf("Error creating bot: %v", err))
	}

	err = bot.Start()
	if err != nil {
		panic(fmt.Sprintf("Error starting bot: %v", err))
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	bot.Close()
}
