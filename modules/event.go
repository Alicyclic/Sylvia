package modules

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Event struct {
	Handler interface{}
}

var Events = make([]*Event, 0)

func AddEvent(handler interface{}) {
	e := &Event{handler}
	Events = append(Events, e)
}

func AddListeners(s *discordgo.Session) {
	for _, event := range Events {
		go func(event *Event) {
			s.AddHandler(event.Handler)
		}(event)
	}
	fmt.Print("Added listeners for events")
}
