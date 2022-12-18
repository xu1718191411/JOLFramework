package framework

import (
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
	h.addHandler("GET", url, combineMiddlewareAndHandler(h.middlewares, handler))
}

func (h *Router) Post(url string, handler func(ctx *JolContext)) {
	h.addHandler("POST", url, combineMiddlewareAndHandler(h.middlewares, handler))
}

func (h *Router) Put(url string, handler func(ctx *JolContext)) {
	h.addHandler("PUT", url, combineMiddlewareAndHandler(h.middlewares, handler))
}

func (h *Router) Patch(url string, handler func(ctx *JolContext)) {
	h.addHandler("PATH", url, combineMiddlewareAndHandler(h.middlewares, handler))
}

func combineMiddlewareAndHandler(middlewares []func(ctx *JolContext), handler func(ctx *JolContext)) []func(ctx *JolContext) {
	arr := make([]func(ctx *JolContext), len(middlewares))
	copy(arr, middlewares)
	return append(arr, handler)
}

func (h *Router) Use(name string, handler func(ctx *JolContext)) {
	if h.middlewares == nil {
		{
			h.middlewares = []func(ctx *JolContext){}
		}
	}

	h.middlewares = append(h.middlewares, handler)
}

func (h *Router) addHandler(method string, url string, handlers []func(ctx *JolContext)) {
	tree := h.handlers[method]
	if tree == nil {
		tree = &Tree{}
		h.handlers[method] = tree
	}
	h.handlers[method].Add(url, handlers)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	successCh := make(chan any)

	jolContext := NewContext(w, r)

	s := e.Router.handlers[strings.ToUpper(r.Method)]
	targetNode := s.Find(r.RequestURI)

	if targetNode == nil {
		jolContext.Status(http.StatusNotFound)
		return
	}

	handlers := targetNode.handlers
	jolContext.Handlers = handlers

	go func() {
		time.Sleep(time.Second)
		jolContext.Next()
		successCh <- 1
	}()

	select {
	case <-successCh:
		fmt.Println("success")
	}

}
