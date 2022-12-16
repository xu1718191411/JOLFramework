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

	tree.Add("/users/:id/tickets", func(ctx *JolContext) {
		ctx.Json("get users id")
	})

	tree.Add("/users/:id/tickets", func(ctx *JolContext) {
		ctx.Json("get users id")
	})

	res := tree.Find("/users/1/lists")

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
