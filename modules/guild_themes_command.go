package modules

import (
	"fmt"
	"time"

	"github.com/carabelle/alexisbot/utils"
	"github.com/thoas/go-funk"
)

// *** GUILD SCHEMES *** //
type GuildTheme struct {
	GuildID, Icon, Name, Banner string
}

var GuildMapping = map[string]map[string]*GuildTheme{
	"966465218127495219": {
		"Christmas": {
			Name: "Celestial Halls o' Jolly!",
		},
		"Halloween": {
			Name: "The Scary Celestial Cafe",
		},
		"Thanksgiving": {
			Name: "Mmm! Celestial Falls!",
		},
		"New Year's Day": {
			Name: "Celestial Cafe",
		},
		"Birthday": {
			Name: "Yummy! Celestial Cake!",
		},
	},
}

//*** HOLIDAYS ***//
type Holiday struct {
	Name  string `json:"name"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
}

type Holidays []*Holiday

var holl = []*Holiday{
	{
		Name:  "New Year's Day",
		Start: 1672531200,
		End:   1675036800,
	},
	{
		Name:  "My Master's Birthday",
		Start: 1672012800,
		End:   1672099200,
	},
	{
		Name:  "Thanksgiving Day",
		Start: 1667260800,
		End:   1669766400,
	},
	{
		Name:  "Christmas",
		Start: 1669852800,
		End:   1672012800,
	},
	{
		Name:  "Halloween",
		Start: 1667174400,
		End:   1667260800,
	},
}

func ConvertDateToEpoch(dateFormat string) int64 {
	t, e := time.Parse("01/02/2006", dateFormat)
	utils.CheckIfError(e)
	return t.AddDate(0, 0, 1).Unix()
}

func ConvertEpochToDate(epoch int64) string {
	t := time.Unix(epoch, 0)
	return t.Format("01/02/2006")
}

func ScheduleForNextYear(e int64) int64 {
	t := time.Unix(e, 0)
	t = t.AddDate(1, 0, 0)
	return t.Unix()
}

func GetHolidayByName(name string) *Holiday {
	if funk.Contains(holl, name) {
		return funk.Get(holl, name).(*Holiday)
	}
	return &Holiday{Name: "No Holiday"}
}

func GetCurrentHoliday() *Holiday {
	isHoliday := funk.Filter(holl, func(h *Holiday) bool {
		return IsHoliday(h)
	})
	if !funk.IsEmpty(isHoliday) {
		return isHoliday.([]*Holiday)[0]
	}
	return &Holiday{Name: "No Holiday"}
}

func IsHoliday(h *Holiday) bool {
	if !funk.IsEmpty(holl) {
		return funk.Filter(holl, func(h *Holiday) bool {
			return h.Start <= time.Now().Unix() && h.End >= time.Now().Unix()
		}).([]*Holiday)[0] == h
	}
	return false
}

func Print() (s string) {
	for _, h := range holl {
		s += fmt.Sprintf("%v-%v: %s\n", ConvertEpochToDate(h.Start), ConvertEpochToDate(h.End), h.Name)
	}
	return s
}
