package cmdlr2

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"sort"
	"strings"
)

type Router struct {
	Prefixes         []string
	IgnorePrefixCase bool
	BotsAllowed      bool
	Commands         []*Command
	Client           *disgord.Client
	Middlewares      []Middleware
	PingHandler      ExecutionHandler
	Storage          map[string]*ObjectsMap
}

func Create(router *Router) *Router {
	router.Storage = map[string]*ObjectsMap{}
	return router
}

func (r *Router) RegisterCMD(command *Command) {
	r.Commands = append(r.Commands, command)
}

func (r *Router) RegisterCMDList(commands []*Command) {
	r.Commands = append(r.Commands, commands...)
}

func (r *Router) GetCmd(name string) *Command {
	sort.Slice(r.Commands, func(i, j int) bool {
		return len(r.Commands[i].Name) > len(r.Commands[j].Name)
	})

	for _, cmd := range r.Commands {
		toCheck := make([]string, len(cmd.Aliases)+1)
		toCheck = append(toCheck, cmd.Name)
		toCheck = append(toCheck, cmd.Aliases...)
		sort.Slice(toCheck, func(i, j int) bool {
			return len(toCheck[i]) > len(toCheck[j])
		})

		if StringArrayContains(toCheck, name, cmd.IgnoreCase) {
			return cmd
		}
	}
	return nil
}

func (r *Router) RegisterMiddleware(middleware Middleware) {
	r.Middlewares = append(r.Middlewares, middleware)
}

func (r *Router) InitializeStorage(name string) {
	r.Storage[name] = NewObjectsMap()
}

func (r *Router) Initialize(client *disgord.Client) {
	client.Gateway().MessageCreate(r.Handler(client))
}

func (r *Router) Handler(c *disgord.Client) (disgord.HandlerMessageCreate, disgord.HandlerMessageCreate) {
	return func(s disgord.Session, h *disgord.MessageCreate) {
		msg := h.Message
		content := h.Message.Content

		if msg.Author.Bot && !r.BotsAllowed {
			return
		}

		u, _ := s.CurrentUser().Get()

		if (content == fmt.Sprintf("<@!%v>", u.ID) || content == fmt.Sprintf("<@%v>", u.ID)) && r.PingHandler != nil {
			r.PingHandler(&Ctx{
				Session: &s,
				Event:   h,
				Client:  r.Client,
				Args:    ParseArguments(""),
				Router:  r,
			})
			return
		}

		hasPrefix, content := StringHasPrefix(content, r.Prefixes, r.IgnorePrefixCase)
		if !hasPrefix {
			return
		}

		content = strings.Trim(content, " ")
		if content == "" {
			return
		}

		for _, m := range r.Middlewares {
			m.Trigger(Ctx{
				Session: &s,
				Event:   h,
				Client:  c,
				Router:  r,
			})
		}

		for _, cmd := range r.Commands {
			toCheck := BuildCheckPrefixes(cmd)

			isCommand, content := StringHasPrefix(content, toCheck, cmd.IgnoreCase)

			if !isCommand {
				continue
			}

			isValid, content := StringHasPrefix(content, []string{" ", "\n"}, false)
			if content == "" || isValid {
				ctx := &Ctx{
					Session: &s,
					Event:   h,
					Args:    ParseArguments(content),
					Router:  r,
					Command: cmd,
				}

				cmd.Trigger(ctx)
			}
		}

	}, nil
}
