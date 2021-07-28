package main

type Stack struct {
	s   [12]addr
	len int
}

func (s *Stack) Len() int    { return s.len }
func (s *Stack) Peek() addr  { return s.s[s.len-1] }
func (s *Stack) Pop() addr   { s.len--; return s.s[s.len] }
func (s *Stack) Push(v addr) { s.s[s.len] = v; s.len++ }
