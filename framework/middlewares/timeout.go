package middlewares

import (
	"JOLFramework/framework"
	"context"
	"time"
)

func Timeout(duration time.Duration) func(ctx *framework.JolContext) {

	return func(ctx *framework.JolContext) {
		c, cancel := context.WithTimeout(ctx, duration)
		defer cancel()

		successCh := make(chan struct{})

		go func(successCh chan struct{}) {
			ctx.Next()
			successCh <- struct{}{}
		}(successCh)

		select {
		case <-c.Done():
			ctx.Lock()
			defer ctx.UnLock()
			ctx.Json("timeout...")
			ctx.SetIsTimeout(true)
		case <-successCh:
			return
		}

	}
}
