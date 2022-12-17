package main

import (
	"JOLFramework/framework"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := framework.NewHandler()

	router.Get("/", func(ctx *framework.JolContext) {
		ctx.Json("hello")
	})

	router.Use("log", func(ctx *framework.JolContext) {
		fmt.Println(ctx.BaseContext())
		ctx.Next()
	})

	group := framework.NewGroup(router, "/api/v1")

	group.Get("/users", func(ctx *framework.JolContext) {
		ctx.Json("users")
	})

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
