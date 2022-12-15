package framework

import (
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

	res := tree.Find("/api/v3/messages")

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

func (m *MockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {

}
