package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/carabelle/alexisbot/utils"
	"github.com/thoas/go-funk"
)

type Command struct {
	*discordgo.Session
	*discordgo.Interaction
	Options map[string]*discordgo.ApplicationCommandInteractionDataOption
}

func (c *Command) CheckIfNSFW() bool {
	channel, _ := c.State.Channel(c.ChannelID)
	return channel.NSFW
}

func (c *Command) CheckIfError(e error) {
	utils.CheckIfError(e)
	if e != nil {
		c.SendErrorMessage(fmt.Sprint(e))
	}
}

func (c *Command) CheckRolePositions(member *discordgo.Member, guild *discordgo.Guild, applyRole *discordgo.Role) bool {
	if !funk.IsEmpty(member.Roles) {
		role := funk.Filter(guild.Roles, func(r *discordgo.Role) bool {
			return r.ID == member.Roles[0]
		}).([]*discordgo.Role)[0]
		return role.Position < applyRole.Position
	}
	return true
}

func (c *Command) HasRole(user *discordgo.User, role *discordgo.Role) bool {
	member, _ := c.Session.GuildMember(c.GuildID, user.ID)
	return funk.Contains(member.Roles, role.ID)
}

func (c *Command) HasPermissionWith(permission int64, member *discordgo.Member) bool {
	p, _ := c.Session.UserChannelPermissions(member.User.ID, c.ChannelID)
	return p&permission == permission
}

func (c *Command) Executor() *discordgo.Member {
	return c.Interaction.Member
}

func (c *Command) HasPermission(permission int64) bool {
	return c.HasPermissionWith(permission, c.Interaction.Member)
}

func (c *Command) GetOption(option string) (value *discordgo.ApplicationCommandInteractionDataOption, ok bool) {
	value, ok = c.Options[option]
	return
}

func (c *Command) parseOptions() {
	options := c.Interaction.ApplicationCommandData().Options
	c.Options = make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		c.Options[opt.Name] = opt
	}
}

func (c *Command) AddOption(name string, t discordgo.ApplicationCommandOptionType) *Command {
	c.Options[name] = &discordgo.ApplicationCommandInteractionDataOption{Name: name, Type: t}
	return c
}

func (c *Command) SendEphemeralMessage(msg string) {
	c.SendEphemeralMessageEmbed(utils.NewEmbed().SetDescription(msg))
}

func (c *Command) Bot() *discordgo.Member {
	member, _ := c.Session.GuildMember(c.GuildID, c.State.User.ID)
	return member
}

func (c *Command) SendEphemeralMessageEmbed(embed *utils.Embed) {
	p, e := c.Session.UserChannelPermissions(c.State.User.ID, c.ChannelID)
	if e != nil {
		c.SendErrorMessage(fmt.Sprint(e))
	}
	if p&discordgo.PermissionEmbedLinks != 0 {
		c.InteractionRespond(c.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:  1 << 6,
				Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
			},
		})
		return
	}
	if embed.MessageEmbed.Description == "" || len(embed.MessageEmbed.Fields) > 0 {
		for _, field := range embed.MessageEmbed.Fields {
			if embed.MessageEmbed.Description != "" {
				embed.MessageEmbed.Description += "\n"
			}
			embed.MessageEmbed.Description += fmt.Sprintf("**%s**: %s", field.Name, field.Value)
		}
	}
	c.InteractionRespond(c.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: embed.Description,
			Flags:   1 << 6,
		},
	})
}

func (c *Command) SendErrorMessage(msg string) {
	c.SendEphemeralMessageEmbed(utils.NewEmbed().SetTitle("Error").SetColor(0xFF0000).SetDescription(msg))
}

func (c *Command) SendInteractionMessage(msg string) {
	c.SendInteractionMessageEmbed(utils.NewEmbed().SetDescription(msg))
}

func (c *Command) SendInteractionMessageEmbed(embed *utils.Embed) {
	p, e := c.Session.UserChannelPermissions(c.State.User.ID, c.ChannelID)
	if e != nil {
		c.SendErrorMessage(fmt.Sprint(e))
	}
	if p&discordgo.PermissionEmbedLinks != 0 {
		c.InteractionRespond(c.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
			},
		})
		return
	}
	if embed.MessageEmbed.Description == "" || len(embed.MessageEmbed.Fields) > 0 {
		for _, field := range embed.MessageEmbed.Fields {
			if embed.MessageEmbed.Description != "" {
				embed.MessageEmbed.Description += "\n"
			}
			embed.MessageEmbed.Description += fmt.Sprintf("**%s**: %s", field.Name, field.Value)
		}
	}
	c.InteractionRespond(c.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: embed.Description,
		},
	})
}
