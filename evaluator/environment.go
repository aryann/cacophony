package evaluator

import "explang/parser"

type environment struct {
	scopes []scope
}

type scope struct {
	vars map[string]parser.Node
}

func newScope() *scope {
	return &scope{
		vars: make(map[string]parser.Node, 0),
	}
}

func newEnvironment() *environment {
	return &environment{
		scopes: []scope{*newScope()},
	}
}

func (e *environment) PushScope() {
	e.scopes = append(e.scopes, *newScope())
}

func (e *environment) PopScope() {
	e.scopes = e.scopes[:len(e.scopes)-1]
}

func (e *environment) Get(name string) (parser.Node, bool) {
	for i := len(e.scopes) - 1; i >= 0; i-- {
		result, ok := e.scopes[i].vars[name]
		if ok {
			return result, true
		}
	}
	return nil, false
}

func (e *environment) Set(name string, value parser.Node) {
	e.scopes[len(e.scopes)-1].vars[name] = value
}
