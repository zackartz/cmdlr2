package cmdlr2

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

func (r *Router) RegisterDefaultHelpCommand(c *disgord.Client) {
	r.InitializeStorage("hdl_helpMessages")

	c.Gateway().MessageReactionAdd(func(s disgord.Session, h *disgord.MessageReactionAdd) {
		channelID := h.ChannelID
		messageID := h.MessageID
		userID := h.UserID
		u, _ := c.Cache().GetCurrentUser()

		if userID == u.ID {
			return
		}

		rawPage, ok := r.Storage["hdl_helpMessages"].Get(fmt.Sprintf("%v:%v:%v", channelID, messageID, userID))
		if !ok {
			return
		}

		page := rawPage.(int)
		if page <= 0 {
			return
		}

		reactionName := h.PartialEmoji.Name
		switch reactionName {
		case "⬅":
			embed, newPage := renderDefaultGeneralHelpEmbed(r, page-1)
			page = newPage
			b := c.Channel(channelID).Message(messageID)
			b.SetEmbed(embed)

			q := c.Channel(channelID).Message(messageID).Reaction(reactionName)
			_ = q.DeleteUser(userID)
			break
		case "❌":
			_ = c.Channel(channelID).Message(messageID).Delete()
			break
		case "➡":
			embed, newPage := renderDefaultGeneralHelpEmbed(r, page+1)
			page = newPage
			b := c.Channel(channelID).Message(messageID)
			b.SetEmbed(embed)

			q := c.Channel(channelID).Message(messageID).Reaction(reactionName)
			_ = q.DeleteUser(userID)
		}

		r.Storage["hdl_helpMessages"].Set(fmt.Sprintf("%v:%v:%v", channelID, messageID, userID), page)
	})

	r.RegisterCMD(&Command{
		Name:        "help",
		Description: "Lists all the available commands or displays some information about a specific command",
		Usage:       "help [command name]",
		Example:     "help yourCommand",
		IgnoreCase:  true,
		Handler:     generalHelpCommand,
	})
}

func generalHelpCommand(ctx *Ctx) {
	if ctx.Args.Amount() > 0 {
		specificHelpCommand(ctx)
		return
	}

	channelID := ctx.Event.Message.ChannelID
	c := ctx.Client

	embed, _ := renderDefaultGeneralHelpEmbed(ctx.Router, 1)
	message, _ := ctx.Client.Channel(channelID).CreateMessage(&disgord.CreateMessageParams{Embed: embed})

	one := c.Channel(channelID).Message(message.ID).Reaction("⬅")
	_ = one.Create()
	two := c.Channel(channelID).Message(message.ID).Reaction("❌")
	_ = two.Create()
	three := c.Channel(channelID).Message(message.ID).Reaction("➡")
	_ = three.Create()

	ctx.Router.Storage["hdl_helpMessages"].Set(fmt.Sprintf("%v:%v:%v", channelID, message.ID, ctx.Event.Message.Author.ID), 1)
}

func specificHelpCommand(ctx *Ctx) {
	commandNames := strings.Split(ctx.Args.Raw(), " ")

	var command *Command
	for index, commandName := range commandNames {
		if index == 0 {
			command = ctx.Router.GetCmd(commandName)
			continue
		}
		command = command.GetSubCommand(commandName)
	}

	_, _ = ctx.Client.Channel(ctx.Event.Message.ChannelID).CreateMessage(&disgord.CreateMessageParams{Embed: renderDefaultSpecificHelpEmbed(ctx, command)})
}

func renderDefaultSpecificHelpEmbed(ctx *Ctx, command *Command) *disgord.Embed {
	prefix := ctx.Router.Prefixes[0]

	if command == nil {
		return &disgord.Embed{
			Type:  "rich",
			Title: "Error",
			Timestamp: disgord.Time{
				Time: time.Now(),
			},
			Color: 0xff0000,
			Fields: []*disgord.EmbedField{
				{
					Name:   "Message",
					Value:  "```The given command doesn't exist. Type `" + prefix + "help` for a list of available commands.```",
					Inline: false,
				},
			},
		}
	}

	subCommands := "No sub commands"
	if len(command.SubCommands) > 0 {
		subCommandNames := make([]string, len(command.SubCommands))
		for index, subCommand := range command.SubCommands {
			subCommandNames[index] = subCommand.Name
		}
		subCommands = "`" + strings.Join(subCommandNames, "`, `") + "`"
	}

	aliases := "No aliases"
	if len(command.Aliases) > 0 {
		aliases = "`" + strings.Join(command.Aliases, "`, `") + "`"
	}

	return &disgord.Embed{
		Title:       "Command Information",
		Type:        "rich",
		Description: "Displaying the information for thr `" + command.Name + "` command.",
		Timestamp: disgord.Time{
			Time: time.Now(),
		},
		Color: 0xffff00,
		Fields: []*disgord.EmbedField{
			{
				Name:   "Name",
				Value:  "`" + command.Name + "`",
				Inline: false,
			},
			{
				Name:   "Sub Commands",
				Value:  subCommands,
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  aliases,
				Inline: false,
			},
			{
				Name:   "Description",
				Value:  "```" + command.Description + "```",
				Inline: false,
			},
			{
				Name:   "Usage",
				Value:  "```" + prefix + command.Usage + "```",
				Inline: false,
			},
			{
				Name:   "Example",
				Value:  "```" + prefix + command.Example + "```",
				Inline: false,
			},
		},
	}
}

func renderDefaultGeneralHelpEmbed(r *Router, page int) (*disgord.Embed, int) {
	commands := r.Commands
	prefix := r.Prefixes[0]

	pageAmount := int(math.Ceil(float64(len(commands)) / 5))
	if page > pageAmount {
		page = pageAmount
	}
	if page <= 0 {
		page = 1
	}

	startingIndex := (page - 1) * 5
	endingIndex := startingIndex + 5
	if page == pageAmount {
		endingIndex = len(commands)
	}
	displayCommands := commands[startingIndex:endingIndex]

	fields := make([]*disgord.EmbedField, len(displayCommands))
	for index, command := range displayCommands {
		fields[index] = &disgord.EmbedField{
			Name:   command.Name,
			Value:  "`" + command.Description + "`",
			Inline: false,
		}
	}

	return &disgord.Embed{
		Title:       fmt.Sprintf("Command List (Page %s / %s)", strconv.Itoa(page), strconv.Itoa(pageAmount)),
		Type:        "rich",
		Description: fmt.Sprintf("These are all the available commands. Type `%shelp <command_name> to find out more about a specific command.`", prefix),
		Timestamp: disgord.Time{
			Time: time.Now(),
		},
		Color:  0xffff00,
		Fields: fields,
	}, page
}
