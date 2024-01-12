package command

import (
	"meet/internal/pkg/tgbot/command"
	"meet/internal/pkg/tgbot/router"
)

type configurator struct {
	startCommand command.Command
}

func NewConfigurator(startCommand command.Command) router.RouterConfigurator {
	return &configurator{startCommand}
}

func (c *configurator) Configure(r router.Router) {
	r.HandleFunc("start", c.startCommand.Handle)
}
