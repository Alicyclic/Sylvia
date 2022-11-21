package modules

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/carabelle/alexisbot/commands"

	"github.com/carabelle/alexisbot/utils"
)

func init() {
	commands.NewCommand("embed", "Send an embed to a channel").
		AddOption("json", "parse json string!", 3, true).
		SetHandler(func(c *commands.Command) {
			arg, _ := c.GetOption("json")
			var embed utils.Embed
			err := json.Unmarshal([]byte(arg.StringValue()), &embed)
			c.CheckIfError(err)
			if !c.HasPermission(discordgo.PermissionEmbedLinks) {
				c.SendInteractionMessageEmbed(&embed)
				return
			}
			c.SendEphemeralMessageEmbed(&embed)
		})
}
