# cmdlr2

cmdlr2 is a command handler framework for the discord wrapper [disgord](https://github.com/andersfylling/disgord). It is heavily inspired (some files are the same) to a framework for discordgo called [dgc](https://github.com/lus/dgc/). I in no way claim to be the original creator of this, this is only a fork for disgord. 

Example starter app:

```go
package main

import (
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"github.com/zackartz/cmdlr2"
	"os"
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.DebugLevel,
}

func main() {
	// Set up a new Disgord client
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_TOKEN"),
		Logger:   log,
	})
	defer client.Gateway().StayConnectedUntilInterrupted()

	router := cmdlr2.Create(&cmdlr2.Router{
		Prefixes:         []string{"$"},
		Client:           client,
		BotsAllowed:      false,
		IgnorePrefixCase: true,
	})

	router.RegisterCMD(&cmdlr2.Command{
		Name:        "ping",
		Description: "It pings.. and yknow.. pongs",
		Usage:       "ping",
		Example:     "ping",
		Handler: func(ctx *cmdlr2.Ctx) {
			ctx.ResponseText("pong")
		},
	})

	router.RegisterDefaultHelpCommand(client)

	router.Initialize(client)
}
```

In this case make sure to set the DISCORD_TOKEN environment variable to the value of your discord token. 

### CONSIDER THIS BETA SOFTWARE

I have a bot that it is using this and it is working however, you may have problems with dgc's implementation of middleware.

Good luck.
