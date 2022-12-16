package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Engine struct {
	Handler *Handler
}

type Handler struct {
	handlers map[string]*Tree
}

func NewHandler() *Handler {
	handlers := make(map[string]*Tree, 0)
	return &Handler{
		handlers: handlers,
	}
}

func (h *Handler) GET(url string, handler func(ctx *JolContext)) {
	h.addHandler("GET", url, handler)
}

func (h *Handler) POST(url string, handler func(ctx *JolContext)) {
	h.addHandler("POST", url, handler)
}

func (h *Handler) PUT(url string, handler func(ctx *JolContext)) {
	h.addHandler("PUT", url, handler)
}

func (h *Handler) PATCH(url string, handler func(ctx *JolContext)) {
	h.addHandler("PATH", url, handler)
}

func (h *Handler) addHandler(method string, url string, handler func(ctx *JolContext)) {
	tree := h.handlers[method]
	if tree == nil {
		tree = &Tree{}
		h.handlers[method] = tree
	}
	h.handlers[method].Add(url, handler)
}

func (h *Handler) HEAD(url string, handler func(ctx *JolContext)) {
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

	s := e.Handler.handlers[strings.ToUpper(r.Method)]
	targetHandler := s.Find(r.RequestURI)

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicCh <- p
			}
		}()
		time.Sleep(time.Second)
		targetHandler(jolContext)
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

func ControllerWithContext(ctx *JolContext) {
	output := struct {
		Msg string `json:"msg"`
	}{
		Msg: "ok",
	}

	ctx.Json(output)
}

func ControllerWithOutContext(w http.ResponseWriter, r *http.Request) {
	output := struct {
		Msg string `json:"msg"`
	}{
		Msg: "ok",
	}

	data, err := json.Marshal(output)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
