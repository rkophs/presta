package vm

import (
	"github.com/rkophs/presta/system"
)

type Stack struct {
	bp     int
	sp     int
	stack  []system.StackEntry
	frames []int
}

func NewStack() *Stack {
	return &Stack{
		bp:     0,
		sp:     0,
		stack:  []system.StackEntry{},
		frames: []int{0},
	}
}

func (s *Stack) Push(entry system.StackEntry) {
	s.stack = append(s.stack, entry)
	s.sp++
}

func (s *Stack) Pop() system.StackEntry {
	s.sp--
	ret := s.stack[s.sp]
	s.stack = s.stack[:s.sp]
	return ret
}

func (s *Stack) Set(offset int, entry system.StackEntry) {
	s.stack[s.bp+offset] = entry
}

func (s *Stack) Fetch(offset int) system.StackEntry {
	return s.stack[s.bp+offset]
}

func (s *Stack) PopFrame() {
	s.stack = s.stack[:s.bp]
	s.sp = s.bp

	frame_len := len(s.frames) - 1
	s.bp = s.frames[frame_len]
	s.frames = s.frames[:frame_len]
}

func (s *Stack) PushFrame() {
	s.frames = append(s.frames, s.bp)
	s.bp = s.sp
}
