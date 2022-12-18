package framework

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type JolContext struct {
	writer    http.ResponseWriter
	request   *http.Request
	mutax     sync.Mutex
	isTimeout bool
	index     int
	Handlers  []func(ctx *JolContext)
}

func NewContext(writer http.ResponseWriter, request *http.Request) *JolContext {
	return &JolContext{
		writer:  writer,
		request: request,
		mutax:   sync.Mutex{},
	}
}

func (c *JolContext) BaseContext() context.Context {
	return c.request.Context()
}

func (c *JolContext) Done() <-chan struct{} {
	return c.BaseContext().Done()
}

func (c *JolContext) Deadline() (deadline time.Time, ok bool) {
	return c.BaseContext().Deadline()
}

func (c *JolContext) Err() error {
	return c.BaseContext().Err()
}

func (c *JolContext) Value(key any) any {
	return c.BaseContext().Value(key)
}

func (c *JolContext) Lock() {
	c.mutax.Lock()
}

func (c *JolContext) UnLock() {
	c.mutax.Unlock()
}

func (c *JolContext) SetIsTimeout(v bool) {
	c.isTimeout = v
}

func (c *JolContext) Send(text string) {
	if c.isTimeout {
		return
	}

	c.writer.WriteHeader(http.StatusOK)
	c.writer.Header().Set("Content-Type", "text/plain")
	c.writer.Write([]byte(text))
}

func (c *JolContext) Json(data any) {
	if c.isTimeout {
		return
	}
	byteData, err := json.Marshal(data)

	if err != nil {
		c.writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.writer.Header().Set("Content-Type", "application/json")
	c.writer.WriteHeader(http.StatusOK)
	c.writer.Write(byteData)
}

func (c *JolContext) Status(code int) {
	if c.isTimeout {
		return
	}
	c.writer.WriteHeader(code)
}

func (c *JolContext) Next() {

	for c.index < len(c.Handlers) {
		handler := c.Handlers[c.index]
		c.index += 1
		handler(c)
	}
}
