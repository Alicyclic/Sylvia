package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/carabelle/alexisbot/utils"
)

type CommandManager struct {
	*discordgo.Session
}

func NewCommandManager(s *discordgo.Session) (b *CommandManager) {
	b = &CommandManager{s}
	b.Init()
	return
}

func (b *CommandManager) RegisterCommandWithin(command *CommandDispatcher, guildId string) (err error) {
	command.GuildId = guildId
	err = b.RegisterCommand(command)
	return
}

func (ds *CommandManager) RegisterCommandsWithin(guildId string, commands []*CommandDispatcher) {
	for _, command := range commands {
		go func(command *CommandDispatcher) {
			err := ds.RegisterCommandWithin(command, guildId)
			utils.CheckIfError(err)
		}(command)
	}
}

func (ds *CommandManager) RegisterCommands(commands []*CommandDispatcher) {
	for _, command := range commands {
		go func(command *CommandDispatcher) {
			err := ds.RegisterCommand(command)
			utils.CheckIfError(err)
		}(command)
	}
}

func (b *CommandManager) RegisterCommand(app *CommandDispatcher) (err error) {
	if CommandExists(app.Name()) {
		log.Printf("Command %s already exists", app.Name())
		return
	}
	if !app.Global() {
		log.Printf("Command %s is guild-specific", app.Name())
	}
	app.Specification, err = b.ApplicationCommandCreate(b.State.User.ID, app.GuildId, app.Specification)
	utils.CheckIfError(err)
	AddCommand(app)
	return nil
}

func (ds *CommandManager) UnregisterCommands() {
	log.Println("Unregistering commands")
	for _, command := range Commands {
		go func(command *CommandDispatcher) {
			err := ds.ApplicationCommandDelete(ds.State.User.ID, command.GuildId, command.Specification.ID)
			utils.CheckIfError(err)
		}(command)
	}
}

func (b *CommandManager) Init() {
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if command, ok := DispatchCommand(i.ApplicationCommandData().Name); ok {
			go command.Invoke(&Command{Session: s, Interaction: i.Interaction})
		}
	})
}
