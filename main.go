package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/carabelle/alexisbot/commands"
	"github.com/carabelle/alexisbot/modules"
)

var (
	commandManager *commands.CommandManager
)

func main() {
	b, _ := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN")))
	if err := b.Open(); err != nil {
		panic(err)
	}
	commandManager = commands.NewCommandManager(b)
	commandManager.RegisterCommandsWithin("989273813055320164", commands.GetCommands())
	// modules.AddListeners(b)
	defer func(session *discordgo.Session) {
		commandManager.UnregisterCommands()
		modules.RemoveAllEphemeralsOnShutdown(session)
		session.Close()
	}(b)
	sc := make(chan os.Signal, 1)
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	signal.Notify(sc, os.Interrupt)
	<-sc
	log.Println("Bot is now closing.")
}
