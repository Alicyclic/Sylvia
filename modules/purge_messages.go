package modules

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/carabelle/alexisbot/commands"
)

func init() {
	commands.NewCommand("purge", "Purge messages from a channel!").
		AddOption("amount", "The amount of messages to purge", 4, true).
		SetHandler(func(c *commands.Command) {
			if !c.HasPermission(discordgo.PermissionManageMessages) {
				c.SendErrorMessage("You don't have permission to manage messages")
			}
			arg, _ := c.GetOption("amount")
			amount := arg.IntValue()
			if amount < 1 || amount >= 100 {
				c.SendErrorMessage("Amount must be between 1 and 100")
				return
			}
			messages, err := c.ChannelMessages(c.ChannelID, int(amount), "", "", "")
			c.CheckIfError(err)
			for _, message := range messages {
				go func(message *discordgo.Message) {
					c.ChannelMessageDelete(c.ChannelID, message.ID)
				}(message)
			}
			c.SendEphemeralMessage(fmt.Sprintf("Purged %d messages", amount))
		})
}
