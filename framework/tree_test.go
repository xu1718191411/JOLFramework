package framework

import (
	"fmt"
	"net/http"
	"testing"
)

func TestTreeT(t *testing.T) {
	tree := &Tree{
		Node: nil,
	}

	tree.Add("/api/v2/messages", func(ctx *JolContext) {
		ctx.Json("get v2 messages")
	})

	tree.Add("/api/v3/messages", func(ctx *JolContext) {
		ctx.Json("get v3 messages")
	})
	tree.Add("/api/v1/tickets", func(ctx *JolContext) {
		ctx.Json("get v1 tickets")
	})
	tree.Add("/users/:id", func(ctx *JolContext) {
		ctx.Json("get users id")
	})
	tree.Add("/", func(ctx *JolContext) {
		ctx.Json("get root")
	})

	res := tree.Find("/api/v1/ttt")

	if res == nil {
		return
	}
	mockResponseWriter := &MockResponseWriter{}
	ctx := NewContext(mockResponseWriter, &http.Request{})
	res(ctx)
}

type MockResponseWriter struct {
}

func (m *MockResponseWriter) Header() http.Header {
	header := make(map[string][]string, 0)
	return header
}

func (m *MockResponseWriter) Write(data []byte) (int, error) {
	fmt.Println(string(data))
	return 0, nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {

}
