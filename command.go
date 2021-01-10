package cmdlr2

import (
	"sort"
	"strings"
)

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Example     string
	Flags       []string
	IgnoreCase  bool
	SubCommands []*Command
	Handler     ExecutionHandler
}

func (c *Command) GetSubCommand(name string) *Command {
	sort.Slice(c.SubCommands, func(i, j int) bool {
		return len(c.SubCommands[i].Name) > len(c.SubCommands[j].Name)
	})

	for _, subCommand := range c.SubCommands {
		toCheck := make([]string, len(subCommand.Aliases)+1)
		toCheck = append(toCheck, subCommand.Name)
		toCheck = append(toCheck, subCommand.Aliases...)
		sort.Slice(toCheck, func(i, j int) bool {
			return len(toCheck[i]) > len(toCheck[j])
		})

		if StringArrayContains(toCheck, name, subCommand.IgnoreCase) {
			return subCommand
		}
	}
	return nil
}

func (c *Command) Trigger(ctx *Ctx) {
	if len(ctx.Args.args) > 0 {
		argument := ctx.Args.Get(0).Raw()
		subCommand := c.GetSubCommand(argument)
		if subCommand != nil {
			args := ParseArguments("")
			if ctx.Args.Amount() > 1 {
				args = ParseArguments(strings.Join(strings.Split(ctx.Args.Raw(), " ")[1:], " "))
			}

			subCommand.Trigger(&Ctx{
				Session: ctx.Session,
				Event:   ctx.Event,
				Args:    args,
				Client:  ctx.Client,
				Router:  ctx.Router,
				Command: subCommand,
			})
			return
		}
	}

	// Prepare all middlewares
	nextHandler := c.Handler

	if nextHandler != nil {
		nextHandler(ctx)
	}
}
