package mc

import (
	"fmt"
	"io"
	"strings"

	"github.com/jbhannah/gophermine/pkg/runner"
)

type CommandType int

const ServerCommands runner.ContextKey = "commands"

const (
	UnknownCommand CommandType = iota
	StopCommand
)

const (
	StopCommandName = "stop"
)

func (cmd CommandType) String() string {
	switch cmd {
	case StopCommand:
		return StopCommandName
	}

	return ""
}

type CommandOrigin interface {
	io.Writer
	Name() string
}

type Command struct {
	CommandType
	Args   []string
	Origin CommandOrigin
}

func NewCommand(origin CommandOrigin, args ...string) *Command {
	return &Command{
		CommandType: stringToCommandType(args[0]),
		Args:        args[1:],
		Origin:      origin,
	}
}

func stringToCommandType(arg string) CommandType {
	switch arg {
	case StopCommandName:
		return StopCommand
	}

	return UnknownCommand
}

func (command *Command) String() string {
	return fmt.Sprintf("[%s] %s %s", command.Origin.Name(), command.CommandType, strings.Join(command.Args, " "))
}
