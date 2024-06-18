package iq_wordbreak

import (
	"sync"
)

type stack struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []*TriNode
}

func NewStack() *stack {
	return &stack{sync.Mutex{}, make([]*TriNode, 0)}
}

func (s *stack) Push(v *TriNode) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}
func (s * stack) Length ()int {
   return len(s.s)
}
func (s *stack) Pop() *TriNode {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return nil
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res
}
