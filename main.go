package main

import (
	"JOLFramework/framework"
	"JOLFramework/framework/middlewares"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	router := framework.NewHandler()

	router.Get("/", func(ctx *framework.JolContext) {
		time.Sleep(time.Second * 3)
		ctx.Json("hello")
	})

	router.Use("recovery", middlewares.Recovery)
	router.Use("log", middlewares.Logger)
	router.Use("timeout", middlewares.Timeout(10*time.Second))

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

	router.Get("/users/:user_id/lists/:list_id", func(ctx *framework.JolContext) {
		ctx.Json("user_id:list_id" + ctx.ParamStringWithDefaultValue(":user_id", "aaa"))
	})

	type User struct {
		Mail     string `json:"mail"`
		Password string `json:"password"`
	}

	router.Post("/posts", func(ctx *framework.JolContext) {
		var user User
		ctx.BindJson(&user)
		var user2 User
		ctx.BindJson(&user2)
		ctx.Json(user)
	})

	engine := framework.Engine{
		Router: router,
	}
	port := 8080

	c := make(chan os.Signal)
	signal.Notify(c)

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: &engine}
	go func(s *http.Server) {
		fmt.Println("listening on port:", port)
		err := s.ListenAndServe()
		if err != nil {
			fmt.Println("error:", err)
		}
	}(server)

	<-c

	fmt.Println("interrupted")

	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Println("error at shutdown", err)
	}
	fmt.Println("shutdown completely")

}
