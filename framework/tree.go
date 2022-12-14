package framework

import (
	"strings"
)

type Node struct {
	param    string
	children []*Node
	handlers []func(ctx *JolContext)
	Parent   *Node
}

func (n *Node) ExistedInChildren(param string) *Node {
	children := n.children

	if len(children) == 0 {
		return nil
	}

	// find normal param is match param or not
	for _, child := range children {
		if child.param == param {
			return child
		}
	}

	// if no noraml param is found, try to find a generic child
	for _, child := range children {
		if isGeneric(child.param) {
			return child
		}
	}

	return nil
}

// is param contains :
func isGeneric(param string) bool {
	return strings.HasPrefix(param, ":")
}

func (n *Node) ParseParams(path string) map[string]string {

	pointer := n
	dicts := make(map[string]string, 0)

	segments := strings.Split(path, "/")

	for i := len(segments) - 1; i > 0; i-- {

		if isGeneric(pointer.param) {
			dicts[pointer.param] = segments[i]
		}

		pointer = pointer.Parent
	}

	return dicts
}

func (n *Node) Find(urls []string) *Node {
	existedNode := n.ExistedInChildren(urls[0])
	if existedNode == nil {
		return nil
	}

	// if it is the last param, judge return value
	if len(urls) == 1 {
		return existedNode
	}

	return existedNode.Find(urls[1:])
}

type Tree struct {
	Node *Node
}

func (t *Tree) Find(url string) *Node {
	arr := strings.Split(url, "/")
	// slash does not exists in the url
	if len(arr) == 1 {
		return nil
	}
	result := t.Node.Find(arr[1:])
	return result
}

func (t *Tree) Add(url string, handlers []func(ctx *JolContext)) *Tree {

	if t.Node == nil {
		t.Node = &Node{
			param: "",
		}
	}

	existedNode := t.Node

	params := strings.Split(url, "/")

	currentNode := existedNode

	for index, param := range params {
		if index == 0 {
			continue
		}

		findInChildren := currentNode.ExistedInChildren(param)

		// if param is not existed in the children
		if findInChildren == nil {
			// add into child
			newChild := &Node{
				param:  param,
				Parent: currentNode,
			}

			// if it is the last node, append handler to the node
			if index == len(params)-1 {
				newChild.handlers = handlers
			}

			// if children is nil, create a new slice holding new generated child
			if currentNode.children == nil {
				currentNode.children = []*Node{
					newChild,
				}
			} else {
				// if children exists, append new generated child into child
				currentNode.children = append(currentNode.children, newChild)
			}

			// move pointer to new generated child
			currentNode = newChild
			continue
		}

		// if param alread existed in child
		currentNode = findInChildren

	}

	return t
}
