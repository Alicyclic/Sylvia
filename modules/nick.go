package modules

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/carabelle/alexisbot/commands"
	"github.com/thoas/go-funk"
)

type Names struct {
	name   string
	emojis []string
}

func (n *Names) String() string {
	var s strings.Builder
	if n.name == "" {
		return ""
	}
	s.WriteString(n.name)
	if len(n.emojis) > 0 {
		funk.ForEach(n.emojis, func(e string) {
			s.WriteString(e)
		})
	}
	return s.String()
}

func init() {
	commands.NewCommand("nick", "You can nickname yourself with this command!").
		AddOption("nickname", "The nickname to set", discordgo.ApplicationCommandOptionString, false).
		SetHandler(func(c *commands.Command) {
			arg, ok := c.GetOption("nickname")
			change := &Names{}
			change.name, change.emojis = arg.StringValue(), []string{"🐲", "🐾⛓"}
			if !ok || change.name == "" {
				change.name = c.Executor().User.Username
			}
			if len(change.name) > 32 {
				c.SendErrorMessage("Name cannot be longer than 32 characters")
			}
			if len(change.String()) > 32 {
				change.emojis = change.emojis[0:1]
				c.SendErrorMessage("Name is too long, using the first badge!")
			}
			c.CheckIfError(c.Session.GuildMemberNickname(c.GuildID, c.Executor().User.ID, change.String()))
			c.SendEphemeralMessage(fmt.Sprintf("Your nickname is %s", change.String()))
		})
}
