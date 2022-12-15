package framework

import (
	"testing"
)

func TestTreeT(t *testing.T) {
	tree := &Tree{
		Node: nil,
	}

	tree.Add("/api/v2/messages")
	tree.Add("/api/v3/messages")
	tree.Add("/api/v1/tickets")
	tree.Add("/users/:id")
}
