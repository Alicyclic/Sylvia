package modules

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/carabelle/alexisbot/commands"
	"github.com/carabelle/alexisbot/utils"
	"github.com/icza/gox/mathx"
)

func times(att1, att2, att3 int64, offset float64) string {
	return fmt.Sprintf("%g", mathx.Round(float64(2.0*att1+2.0*att2+att3)*offset, 0.1)) + "%"
}

func div(att1, att2, att3 int64, offset float64) string {
	return fmt.Sprintf("%g", mathx.Round(float64(2.0*att1+2.0*att2+att3)/offset, 0.1)) + "%"
}

func calculate(strength, intellect, agility, will, power int64) map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"damage": {
			"bringer": div(strength, will, power, 400.0),
			"giver":   div(strength, will, power, 200.0),
			"dealer":  times(strength, will, power, 0.0075),
		},
		"resist": {
			"ward":  times(strength, agility, power, 0.012),
			"proof": div(strength, agility, power, 125.0),
		},
		"critical": {
			"defender": times(intellect, will, power, 0.024),
			"blocker":  times(intellect, will, power, 0.02),
		},
		"pierce": {
			"breaker": div(strength, agility, power, 400.0),
			"piercer": times(strength, agility, power, 0.0015),
		},
		"stun": {
			"recal":  div(strength, intellect, power, 125.0),
			"resist": div(strength, intellect, power, 250.0),
		},
		"healing": {
			"lively": times(strength, agility, power, 0.0065),
			"healer": times(strength, agility, power, 0.003),
			"medic":  times(strength, agility, power, 0.0065),
		},
		"health": {
			"healthy": times(intellect, agility, power, 0.003),
			"gift":    times(agility, will, power, 0.1),
			"add":     times(agility, will, power, 0.06),
		},
		"attributes": {
			"strength":  strength,
			"intellect": intellect,
			"agility":   agility,
			"willpower": will,
			"power":     power,
			"happiness": strength + intellect + agility + will + power,
		},
	}
}

func embedModel(strength, intellect, agility, will, power int64) (e *utils.Embed) {
	e = utils.NewEmbed().SetTitle("Model")
	resp := calculate(strength, intellect, agility, will, power)

	attributes := []string{
		fmt.Sprintf("Strength %d", strength),
		fmt.Sprintf("Intellect %d", intellect),
		fmt.Sprintf("Agility %d", agility),
		fmt.Sprintf("Will %d", will),
		fmt.Sprintf("Power %d", power),
	}

	critical := []string{
		fmt.Sprintf("Blocker %s", resp["critical"]["blocker"].(string)),
		fmt.Sprintf("Defender %s", resp["critical"]["defender"].(string)),
	}

	damage := []string{
		fmt.Sprintf("Bringer %s", resp["damage"]["bringer"].(string)),
		fmt.Sprintf("Dealer %s", resp["damage"]["dealer"].(string)),
		fmt.Sprintf("Giver %s", resp["damage"]["giver"].(string)),
	}

	healing := []string{
		fmt.Sprintf("Healer %s", resp["healing"]["healer"].(string)),
		fmt.Sprintf("Lively %s", resp["healing"]["lively"].(string)),
		fmt.Sprintf("Medic %s", resp["healing"]["medic"].(string)),
	}

	health := []string{
		fmt.Sprintf("Add %s", resp["health"]["add"].(string)),
		fmt.Sprintf("Gift %s", resp["health"]["gift"].(string)),
		fmt.Sprintf("Healthy %s", resp["health"]["healthy"].(string)),
	}

	piercing := []string{
		fmt.Sprintf("Breaker %s", resp["pierce"]["breaker"].(string)),
		fmt.Sprintf("Piercer %s", resp["pierce"]["piercer"].(string)),
	}

	resist := []string{
		fmt.Sprintf("Proof %s", resp["resist"]["proof"].(string)),
		fmt.Sprintf("Ward %s", resp["resist"]["ward"].(string)),
	}

	stun := []string{
		fmt.Sprintf("Recal %s", resp["stun"]["recal"].(string)),
		fmt.Sprintf("Resist %s", resp["stun"]["resist"].(string)),
	}

	e.AddField("Attributes", strings.Join(attributes, "\n"))
	e.AddField("Critical", strings.Join(critical, "\n"))
	e.AddField("Damage", strings.Join(damage, "\n"))
	e.AddField("Healing", strings.Join(healing, "\n"))
	e.AddField("Health", strings.Join(health, "\n"))
	e.AddField("Piercing", strings.Join(piercing, "\n"))
	e.AddField("Resist", strings.Join(resist, "\n"))
	e.AddField("Stun", strings.Join(stun, "\n"))
	return e
}

func init() {
	commands.NewCommand("pet", "Calculate pet attributes").
		AddOption("strength", "Give a strength attribute?", discordgo.ApplicationCommandOptionInteger, true).
		AddOption("intellect", "Give a intellect attribute?", discordgo.ApplicationCommandOptionInteger, true).
		AddOption("agility", "Give a agility attribute?", discordgo.ApplicationCommandOptionInteger, true).
		AddOption("will", "Give a will attribute?", discordgo.ApplicationCommandOptionInteger, true).
		AddOption("power", "Give a power attribute?", discordgo.ApplicationCommandOptionInteger, true).
		SetHandler(func(c *commands.Command) {
			strength, _ := c.GetOption("strength")
			intellect, _ := c.GetOption("intellect")
			agility, _ := c.GetOption("agility")
			will, _ := c.GetOption("will")
			power, _ := c.GetOption("power")

			e := embedModel(strength.IntValue(), intellect.IntValue(), agility.IntValue(), will.IntValue(), power.IntValue())

			c.SendEphemeralMessageEmbed(e)
		})
}
