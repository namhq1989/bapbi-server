package waiter

import (
	"context"
)

type Options func(c *waiterCfg)

func ParentContext(ctx context.Context) Options {
	return func(c *waiterCfg) {
		c.parentCtx = ctx
	}
}

func CatchSignals() Options {
	return func(c *waiterCfg) {
		c.catchSignals = true
	}
}
