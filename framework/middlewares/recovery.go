package middlewares

import "JOLFramework/framework"

func Recovery(ctx *framework.JolContext) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Lock()
			defer ctx.UnLock()
			ctx.Json("something error")
		}
	}()
	ctx.Next()
}
