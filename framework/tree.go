package framework

import "strings"

type Node struct {
	param    string
	children []*Node
}

func (n *Node) ExistedInChildren(param string) *Node {
	children := n.children

	if len(children) == 0 {
		return nil
	}

	for _, child := range children {
		if child.param == param {
			return child
		}
	}

	return nil
}

type Tree struct {
	Node *Node
}

func (t *Tree) Add(url string) *Tree {

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
				param: param,
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
