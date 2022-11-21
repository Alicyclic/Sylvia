package modules

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/carabelle/alexisbot/commands"
	"github.com/carabelle/alexisbot/utils"
)

const (
	EphemeralChannelType = 0
	EphemeralRoleType    = 1
)

var (
	Ephemerals            = make([]*Ephemeral, 0)
	expiresArgument int64 = 15
)

type Ephemeral struct {
	*discordgo.Session
	guildId, userId, itemId  string
	EphemeralType, expiresAt int64
}

func AddEphemeral(e *Ephemeral) {
	Ephemerals = append(Ephemerals, e)
	e.CreateEphemeralInGuild()
	duration := time.Duration(e.expiresAt) * time.Second
	time.AfterFunc(duration, func() {
		e.RemoveEphemeral()
	})
}

func (e *Ephemeral) RemoveEphemeral() {
	switch e.EphemeralType {
	case EphemeralChannelType:
		_, err := e.ChannelDelete(e.itemId)
		utils.CheckIfError(err)
	case EphemeralRoleType:
		err := e.GuildMemberRoleRemove(e.guildId, e.userId, e.itemId)
		utils.CheckIfError(err)
	}
}

func RemoveAllEphemeralsOnShutdown(s *discordgo.Session) {
	for _, v := range Ephemerals {
		go func(v *Ephemeral) {
			v.RemoveEphemeral()
		}(v)
	}
}

func (e *Ephemeral) CreateEphemeralInGuild() {
	switch e.EphemeralType {
	case EphemeralChannelType:
		ch, err := e.GuildChannelCreate(e.guildId, "ephemeral-channel-"+e.userId, EphemeralChannelType)
		e.itemId = ch.ID
		utils.CheckIfError(err)
	case EphemeralRoleType:
		err := e.GuildMemberRoleAdd(e.guildId, e.userId, e.itemId)
		utils.CheckIfError(err)
	}
}

func (e *Ephemeral) String() string {
	if e.EphemeralType == EphemeralChannelType {
		return "Ephemeral channel"
	}
	return "Ephemeral role"
}

func (e *Ephemeral) EmbedModel() (m *utils.Embed) {
	m = utils.NewEmbed()
	m.SetTitle("Ephemeral")
	m.SetDescription(fmt.Sprintf("%s will be removed in %d seconds", e.String(), e.expiresAt))
	m.SetColor(0x00FF00)
	return m
}

func init() {
	commands.NewCommand("role", "Give or create a ephemeral role! (wip)").
		AddOption("role", "Give a role", discordgo.ApplicationCommandOptionRole, true).
		AddOption("user", "Give a user", discordgo.ApplicationCommandOptionUser, true).
		AddOption("expires", "How long the role should last", discordgo.ApplicationCommandOptionInteger, true).
		SetHandler(func(c *commands.Command) {
			roleParam, _ := c.GetOption("role")
			userParam, _ := c.GetOption("user")
			expiresParam, ok := c.GetOption("expires")
			if ok {
				expiresArgument = expiresParam.IntValue()
			}
			role, user := roleParam.RoleValue(c.Session, c.GuildID), userParam.UserValue(c.Session)
			guild, _ := c.Session.Guild(c.GuildID)
			bot, _ := c.GuildMember(c.GuildID, c.Session.State.User.ID)
			switch {
			case c.HasRole(user, role):
				c.SendErrorMessage(fmt.Sprintf("%s already has the role %s", user.Mention(), role.Mention()))
				return
			case c.CheckRolePositions(c.Executor(), guild, role):
				c.SendErrorMessage("You can't give a role higher than your own")
			case c.CheckRolePositions(bot, guild, role):
				c.SendErrorMessage("I can't give a role higher than my own")
				return
			case !c.HasPermissionWith(discordgo.PermissionManageRoles, bot):
				c.SendErrorMessage("I don't have permission to manage roles")
				return
			case !c.HasPermission(discordgo.PermissionManageRoles):
				c.SendErrorMessage("You don't have permission to manage roles")
				return
			}
			eph := &Ephemeral{
				Session:       c.Session,
				guildId:       c.GuildID,
				userId:        user.ID,
				itemId:        role.ID,
				EphemeralType: 1,
				expiresAt:     expiresArgument,
			}
			AddEphemeral(eph)
			c.SendEphemeralMessageEmbed(eph.EmbedModel())
		})
	commands.NewCommand("channel", "Give or create a ephemeral channel! (wip)").
		AddOption("expires", "How long the channel should last", discordgo.ApplicationCommandOptionInteger, true).
		SetHandler(func(c *commands.Command) {
			expiresParam, ok := c.GetOption("expires")
			if ok {
				expiresArgument = expiresParam.IntValue()
			}
			bot, _ := c.GuildMember(c.GuildID, c.Session.State.User.ID)
			switch {
			case !c.HasPermissionWith(discordgo.PermissionManageChannels, bot):
				c.SendErrorMessage("I don't have permission to manage channel")
				return
			case !c.HasPermission(discordgo.PermissionManageChannels):
				c.SendErrorMessage("You don't have permission to manage channel")
				return
			}
			eph := &Ephemeral{
				Session:       c.Session,
				guildId:       c.GuildID,
				userId:        c.Executor().User.ID,
				EphemeralType: 0,
				expiresAt:     expiresArgument,
			}
			AddEphemeral(eph)
			c.SendEphemeralMessageEmbed(eph.EmbedModel())
		})
}
