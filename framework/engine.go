package framework

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Engine struct {
	Router *Router
}

type Router struct {
	handlers    map[string]*Tree
	middlewares []func(ctx *JolContext)
}

func (r *Router) GetHandlers() map[string]*Tree {
	return r.handlers
}

func NewHandler() *Router {
	handlers := make(map[string]*Tree, 0)
	handlers["GET"] = &Tree{}
	handlers["POST"] = &Tree{}
	handlers["PUT"] = &Tree{}
	handlers["PATCH"] = &Tree{}
	handlers["DELETE"] = &Tree{}
	handlers["HEAD"] = &Tree{}
	return &Router{
		handlers: handlers,
	}
}

func (h *Router) Get(url string, handler func(ctx *JolContext)) {
	h.addHandler("GET", url, handler)
}

func (h *Router) Post(url string, handler func(ctx *JolContext)) {
	h.addHandler("POST", url, handler)
}

func (h *Router) Put(url string, handler func(ctx *JolContext)) {
	h.addHandler("PUT", url, handler)
}

func (h *Router) Patch(url string, handler func(ctx *JolContext)) {
	h.addHandler("PATH", url, handler)
}

func (h *Router) Use(name string, handler func(ctx *JolContext)) {
	if h.middlewares == nil {
		{
			h.middlewares = []func(ctx *JolContext){}
		}
	}

	h.middlewares = append(h.middlewares, handler)
}

func (h *Router) addHandler(method string, url string, handler func(ctx *JolContext)) {
	tree := h.handlers[method]
	if tree == nil {
		tree = &Tree{}
		h.handlers[method] = tree
	}
	h.handlers[method].Add(url, handler)
}

func (h *Router) HEAD(url string, handler func(ctx *JolContext)) {
	tree := h.handlers["HEAD"]
	if tree == nil {
		tree = &Tree{}
		h.handlers["HEAD"] = tree
	}
	h.handlers["HEAD"].Add(url, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	panicCh := make(chan any)
	successCh := make(chan any)

	jolContext := NewContext(w, r)

	ctx, cancel := context.WithTimeout(jolContext, time.Second*5)

	defer cancel()

	s := e.Router.handlers[strings.ToUpper(r.Method)]
	targetHandler := s.Find(r.RequestURI)
	middlerwares := e.Router.middlewares

	handlers := append(middlerwares, targetHandler)
	jolContext.Handlers = handlers

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicCh <- p
			}
		}()
		time.Sleep(time.Second)
		jolContext.Next()
		successCh <- 1
	}()

	select {
	case <-ctx.Done():
		jolContext.Lock()
		defer jolContext.UnLock()
		jolContext.Send("timeout")
		jolContext.setIsTimeout(true)
	case <-panicCh:
		jolContext.Lock()
		defer jolContext.UnLock()
		jolContext.Send("internal error")
	case <-successCh:
		fmt.Println("success")
	}

}
