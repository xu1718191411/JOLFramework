package main

import (
	"JOLFramework/framework"
	"JOLFramework/framework/middlewares"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	router := framework.NewHandler()

	router.Get("/", func(ctx *framework.JolContext) {
		ctx.Json("hello")
	})

	router.Use("recovery", middlewares.Recovery)
	router.Use("log", middlewares.Logger)
	router.Use("timeout", middlewares.Timeout(time.Second))

	group := framework.NewGroup(router, "/api/v1")

	group.Use("log", func(ctx *framework.JolContext) {
		fmt.Println("Group middleware")
		ctx.Next()
	})

	router.Get("/panic", func(ctx *framework.JolContext) {
		ctx.Json("panic")
	})
	router.Get("/timeout", func(ctx *framework.JolContext) {
		time.Sleep(time.Second * 3)
		ctx.Json("timeout handler")
	})

	group.Get("/tickets", func(ctx *framework.JolContext) {
		ctx.Json("tickets, name:" + ctx.QueryStringWithDefault("name", "defaultName"))
	})

	router.Get("/users", func(ctx *framework.JolContext) {
		ctx.Json(fmt.Sprintf("users, id: %d", ctx.QueryIntWithDefault("user_id", 0)))
	})

	engine := framework.Engine{
		Router: router,
	}
	port := 8080
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), &engine))
	fmt.Println("listening on port:", port)
}
