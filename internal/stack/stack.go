package stack

import (
	"github.com/ethereum/go-ethereum/common"
)

// github.com/golang-collections
type (
	Stack struct {
		top    *node
		length int
	}
	node struct {
		value *common.Address
		prev  *node
	}
)

// Create a new stack
func NewStack() *Stack {
	return &Stack{nil, 0}
}

// Return the number of items in the stack
func (this *Stack) Len() int {
	return this.length
}

// View the top item on the stack
func (this *Stack) Peek() *common.Address {
	if this.length == 0 {
		return nil
	}
	return this.top.value
}

// Pop the top item of the stack and return it
func (this *Stack) Pop() *common.Address {
	if this.length == 0 {
		return nil
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

// Push a value onto the top of the stack
func (this *Stack) Push(value *common.Address) {
	n := &node{value, this.top}
	this.top = n
	this.length++
}
