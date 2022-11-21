package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type CommandDispatcher struct {
	Specification *discordgo.ApplicationCommand
	GuildId       string
	Handle        func(*Command)
}

func (a *CommandDispatcher) Global() bool {
	return a.GuildId == ""
}

func (a *CommandDispatcher) SetGuildId(guildId string) *CommandDispatcher {
	a.GuildId = guildId
	return a
}

func (a *CommandDispatcher) Invoke(cs *Command) {
	defer func() {
		if r := recover(); r != nil {
			cs.SendErrorMessage(fmt.Sprintf("Error: %s", r))
		}
	}()
	cs.parseOptions()
	a.Handle(cs)
}

func (a *CommandDispatcher) Name() string {
	return a.Specification.Name
}

var (
	CommandSlash = make([]*CommandDispatcher, 0)
	Commands     = make(map[string]*CommandDispatcher)
)

func NewCommand(name, desc string) (c *CommandDispatcher) {
	c = &CommandDispatcher{
		Specification: &discordgo.ApplicationCommand{
			Name:        name,
			Description: desc,
		},
	}
	CommandSlash = append(CommandSlash, c)
	return
}

func GetCommands() []*CommandDispatcher {
	return CommandSlash
}

func AddCommand(command *CommandDispatcher) {
	Commands[command.Name()] = command
}

func (c *CommandDispatcher) SetHandler(h func(*Command)) *CommandDispatcher {
	c.Handle = h
	return c
}

func (c *CommandDispatcher) AddOptions(n, d string, t discordgo.ApplicationCommandOptionType, r bool, dc []*discordgo.ApplicationCommandOptionChoice) *CommandDispatcher {
	c.Specification.Options = append(c.Specification.Options, &discordgo.ApplicationCommandOption{
		Name:        n,
		Type:        t,
		Description: d,
		Required:    true,
		Choices:     dc,
	})
	return c
}

func (c *CommandDispatcher) AddOption(n, d string, t discordgo.ApplicationCommandOptionType, r bool) *CommandDispatcher {
	c.Specification.Options = append(c.Specification.Options, &discordgo.ApplicationCommandOption{Name: n, Type: t, Description: d, Required: r})
	return c
}

func DispatchCommand(c string) (command *CommandDispatcher, ok bool) {
	command, ok = Commands[c]
	return
}

func CommandExists(c string) bool {
	_, ok := DispatchCommand(c)
	return ok
}

func (c *CommandDispatcher) SetName(name string) *CommandDispatcher {
	c.Specification.Name = name[:128]
	return c
}

func (c *CommandDispatcher) SetDescription(desc string) *CommandDispatcher {
	c.Specification.Description = desc[:2048]
	return c
}
