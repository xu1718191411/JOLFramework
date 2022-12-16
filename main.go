package main

import (
	"JOLFramework/framework"
	"fmt"
	"log"
	"net/http"
)

func main() {
	handler := framework.NewHandler()

	handler.GET("/api/v1/users", func(ctx *framework.JolContext) {
		ctx.Json("users")
	})

	handler.GET("/", func(ctx *framework.JolContext) {
		ctx.Json("hello")
	})

	handler.GET("/api/v1/users/:id", func(ctx *framework.JolContext) {
		ctx.Json("user detail page")
	})

	engine := framework.Engine{
		Handler: handler,
	}
	port := 8080
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), &engine))
	fmt.Println("listening on port:", port)
}
