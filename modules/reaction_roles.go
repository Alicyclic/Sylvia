package modules

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/carabelle/alexisbot/utils"
	"github.com/thoas/go-funk"
)

type ReactEmoji struct {
	roleId, emojiId, channelId, messageId string
	stackable                             bool
}

type ReactMessage struct {
	emojis []ReactEmoji
}

func (r *ReactMessage) AddRole(roleId, emojiId, channelId, messageId string) {
	r.emojis = append(r.emojis, ReactEmoji{
		roleId:    roleId,
		emojiId:   emojiId,
		channelId: channelId,
		messageId: messageId,
		stackable: true,
	})

	log.Println("Reaction added:", roleId, emojiId, channelId, messageId)
}

func (r *ReactMessage) RemoveRole(roleId, emojiId, channelId, messageId string) {
	for i, emoji := range r.emojis {
		if emoji.roleId == roleId && emoji.emojiId == emojiId && emoji.channelId == channelId && emoji.messageId == messageId {
			r.emojis = append(r.emojis[:i], r.emojis[i+1:]...)
			log.Println("Reaction removed:", roleId, emojiId, channelId, messageId)
			return
		}
	}
}

var react = ReactMessage{}

func init() {
	AddEvent(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.UserID == s.State.User.ID {
			return
		}
		parserEmoji := fmt.Sprintf("%s:%s", m.Emoji.Name, m.Emoji.ID)
		fmt.Println("Reaction:", parserEmoji)
		react.AddRole("987510649930268683", parserEmoji, m.ChannelID, m.MessageID)
		anyRole := funk.Find(react.emojis, func(emoji ReactEmoji) bool {
			return emoji.emojiId == parserEmoji
		}).(ReactEmoji)
		utils.CheckIfError(s.GuildMemberRoleAdd(m.GuildID, m.UserID, anyRole.roleId))
	})

	AddEvent(func(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
		if m.UserID == s.State.User.ID {
			return
		}
		parserEmoji := fmt.Sprintf("%s:%s", m.Emoji.Name, m.Emoji.ID)
		fmt.Println("Reaction:", parserEmoji)
		anyRole := funk.Find(react.emojis, func(emoji ReactEmoji) bool {
			return emoji.emojiId == parserEmoji
		}).(ReactEmoji)
		utils.CheckIfError(s.GuildMemberRoleRemove(m.GuildID, m.UserID, anyRole.roleId))
		log.Println("Reaction removed:", m.UserID, m.MessageID, fmt.Sprintf("%s:%s", m.Emoji.Name, m.Emoji.ID))
	})

}
