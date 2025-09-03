package stack

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// github.com/golang-collections
type (
	Stack struct {
		top    *node
		length int
	}
	node struct {
		value *callAddr
		prev  *node
	}
)

type callAddr struct {
	Addr     *common.Address
	Calltype vm.OpCode
}

// Create a new stack
func NewStack() *Stack {
	return &Stack{nil, 0}
}

// Return the number of items in the stack
func (s *Stack) Len() int {
	return s.length
}

// View the top item on the stack
func (s *Stack) Peek() *common.Address {
	c := s.top

	for c != nil {
		if c.value.Calltype != vm.DELEGATECALL {
			return c.value.Addr
		} else {
			c = c.prev
		}
	}

	return nil
}

// Pop the top item of the stack and return it
func (s *Stack) Pop() *common.Address {
	// Pop 데이터를 별도로 활용하지 않음
	if s.length == 0 {
		return nil
	}

	n := s.top
	s.top = n.prev
	s.length--
	return n.value.Addr
}

// Push a value onto the top of the stack
func (s *Stack) Push(value *common.Address, op vm.OpCode) {
	cAddr := &callAddr{value, op}
	n := &node{cAddr, s.top}
	s.top = n
	s.length++
}
