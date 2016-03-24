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

package parser

type Semantic struct {
	fns   map[string]*fnTuple
	vars  []*scope
	ids   int
	scope int
}

type fnTuple struct {
	arity int
	id    int
}

type scope struct {
	vars map[string]int
}

func NewSemantic() *Semantic {
	return &Semantic{fns: make(map[string]*fnTuple), vars: make([]*scope, 0), ids: 0, scope: 0}
}

func newFnType(arity, id int) *fnTuple {
	return &fnTuple{arity: arity, id: id}
}

func newScope(vars map[string]int) *scope {
	return &scope{vars: vars}
}

func (s *Semantic) nextId() int {
	s.ids++
	return s.ids
}

func (s *Semantic) AddFunction(name string, arity int) {
	s.fns[name] = newFnType(arity, s.nextId())
}

func (s *Semantic) FunctionExists(name string) bool {
	_, ok := s.fns[name]
	return ok
}

func (s *Semantic) FunctionArity(name string) int {
	if fn, ok := s.fns[name]; ok {
		return fn.arity
	} else {
		return -1
	}
}

func (s *Semantic) PushNewScope(names []string) {
	vars := make(map[string]int)
	for _, name := range names {
		vars[name] = s.nextId()
	}

	s.vars = append([]*scope{newScope(vars)}, s.vars...)
}

func (s *Semantic) PopScope() {
	s.vars = s.vars[1:]
}

func (s *Semantic) VariableExists(name string) bool {
	for _, scope := range s.vars {
		if _, ok := scope.vars[name]; ok {
			return true
		}
	}
	return false
}

func (s *Semantic) GetVariableId(name string) int {
	for _, scope := range s.vars {
		if v, ok := scope.vars[name]; ok {
			return v
		}
	}
	return -1
}
