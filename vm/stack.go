/*
 * Copyright (c) 2016 Ryan Kophs
 *
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 **/

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
