package middlewares

import (
	"JOLFramework/framework"
	"fmt"
)

func Logger(ctx *framework.JolContext) {
	fmt.Println("middleware log")
	ctx.Next()
}
