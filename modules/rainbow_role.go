package modules

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/carabelle/alexisbot/utils"
	"github.com/robfig/cron"
	"lukechampine.com/frand"
)

type RainbowRole struct {
	*discordgo.Session
	*discordgo.Ready

	guildId, roleId string
}

func (rainbow *RainbowRole) Run() {
	colors := []string{"ff0004", "ff1493", "ff00ff", "f9a328", "e4ff00", "00ffff", "6363ff"}
	frand.NewSource().Seed(time.Now().Unix())
	colorInt, err := strconv.ParseInt(colors[frand.Intn(len(colors))], 16, 32)
	utils.CheckIfError(err)
	_, err = rainbow.GuildRoleEdit(rainbow.guildId, rainbow.roleId, "Rainbow", int(colorInt), false, 0, false)
	utils.CheckIfError(err)
}

func init() {
	AddEvent(func(s *discordgo.Session, r *discordgo.Ready) {
		RainbowRole := &RainbowRole{s, r, "985809699343572992", "987510649930268683"}
		c := cron.New()
		c.AddFunc("@every 200s", RainbowRole.Run)
		c.Start()
	})

}
