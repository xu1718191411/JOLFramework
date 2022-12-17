package main

import (
	"JOLFramework/framework"
	"JOLFramework/framework/middlewares"
	"fmt"
	"log"
	"net/http"
	"time"
)

func panic(ctx *framework.JolContext) {
	ctx.Json("panic")
}

func timeout(ctx *framework.JolContext) {
	ctx.Json("timeout handler")
}

func main() {
	router := framework.NewHandler()

	// router.Get("/", func(ctx *framework.JolContext) {
	// 	ctx.Json("hello")
	// })

	router.Use("recovery", middlewares.Recovery)
	router.Use("log", middlewares.Logger)
	router.Use("timeout", middlewares.Timeout(time.Second*500))

	group := framework.NewGroup(router, "/api/v1")

	group.Use("log", func(ctx *framework.JolContext) {
		fmt.Println("Group middleware")
		ctx.Next()
	})

	// group.Get("/users", func(ctx *framework.JolContext) {
	// 	ctx.Json("users")
	// })

	router.Get("/panic", panic)
	router.Get("/timeout", timeout)

	group.Get("/tickets", func(ctx *framework.JolContext) {
		ctx.Json("tickets")
	})

	engine := framework.Engine{
		Router: router,
	}
	port := 8080
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), &engine))
	fmt.Println("listening on port:", port)
}
