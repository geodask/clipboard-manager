package cli

import (
	"context"

	"github.com/geodask/clipboard-manager/internal/client"
)

type Command interface {
	Name() string
	Description() string
	Usage() string
	Execute(ctx context.Context, client *client.Client, args []string) error
}

type Registry struct {
	commands map[string]Command
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

func (r *Registry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

func (r *Registry) Get(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

func (r *Registry) All() []Command {
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}
