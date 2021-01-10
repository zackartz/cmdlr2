package cmdlr2

import (
	"context"
	"github.com/andersfylling/disgord"
)

type Ctx struct {
	Client  *disgord.Client
	Session *disgord.Session
	Event   *disgord.MessageCreate
	Args    *Arguments
	Command *Command
	Router  *Router
}

type ExecutionHandler func(ctx *Ctx)

func (ctx *Ctx) ResponseText(text string) error {
	channel, err := ctx.Client.Channel(ctx.Event.Message.ChannelID).Get()
	if err != nil {
		return err
	}
	_, err = channel.SendMsgString(context.Background(), *ctx.Session, text)
	return err
}

func (ctx *Ctx) ResponseEmbed(embed *disgord.Embed) error {
	_, err := ctx.Client.Channel(ctx.Event.Message.ChannelID).CreateMessage(&disgord.CreateMessageParams{
		Embed: embed,
	})
	return err
}

func (ctx *Ctx) ResponseTextEmbed(text string, embed *disgord.Embed) error {
	_, err := ctx.Client.Channel(ctx.Event.Message.ChannelID).CreateMessage(&disgord.CreateMessageParams{
		Embed:   embed,
		Content: text,
	})
	return err
}
