package mc

import (
	"fmt"
	"io"
	"strings"

	"github.com/jbhannah/gophermine/pkg/runner"
)

// ServerCommands is the key of the server command channel in subcontexts of the
// running server.
const ServerCommands runner.ContextKey = "commands"

// CommandType indicates the command being sent.
type CommandType int

const (
	// UnknownCommand means the command sent was not recognized.
	UnknownCommand CommandType = iota

	// StopCommand is a /stop command to stop the server.
	StopCommand
)

// Constants of command keywords.
const (
	StopCommandName = "stop"
)

// String maps a CommandType to its keyword.
func (cmd CommandType) String() string {
	switch cmd {
	case StopCommand:
		return StopCommandName
	}

	return ""
}

// Origin is the input (e.g. RCON, stdin) that sent a given command.
type Origin interface {
	io.Writer
	Name() string
}

// Command represents a command sent to the server.
type Command struct {
	CommandType
	Origin
	Args []string
}

// NewCommand instantiates a Command from the origin and input string.
func NewCommand(origin Origin, args ...string) *Command {
	return &Command{
		CommandType: stringToCommandType(args[0]),
		Origin:      origin,
		Args:        args[1:],
	}
}

func stringToCommandType(arg string) CommandType {
	switch arg {
	case StopCommandName:
		return StopCommand
	}

	return UnknownCommand
}

// String returns a string representation of a command.
func (command *Command) String() string {
	return fmt.Sprintf("[%s] %s %s", command.Origin.Name(), command.CommandType, strings.Join(command.Args, " "))
}
