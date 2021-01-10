package cmdlr2

type Middleware struct {
	Trigger func(ctx Ctx)
}
