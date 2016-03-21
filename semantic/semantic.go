package semantic

import (
	"fmt"
)

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

func newFnTupe(arity, id int) *fnTuple {
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
	s.fns[name] = newFnTupe(arity, s.nextId())
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
		fmt.Println("Scope it")
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
