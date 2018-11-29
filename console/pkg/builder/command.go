package builder

import (
	"hidevops.io/hiboot/pkg/app"
)

type CommandBuilder interface {
}

type Command struct {
	CommandBuilder
}

func init() {
	app.Register(CommandService)
}

func CommandService() CommandBuilder {
	return &Command{}
}

func (c *Command) Send() error {

	return nil
}
