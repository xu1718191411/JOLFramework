package main

import (
	"JOLFramework/framework"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := framework.NewHandler()

	router.GET("/", func(ctx *framework.JolContext) {
		ctx.Json("hello")
	})

	group := router.Group("/api/v1")

	group.GET("/users", func(ctx *framework.JolContext) {
		ctx.Json("users")
	})

	group.GET("/tickets", func(ctx *framework.JolContext) {
		ctx.Json("tickets")
	})

	engine := framework.Engine{
		Router: router,
	}
	port := 8080
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), &engine))
	fmt.Println("listening on port:", port)
}
