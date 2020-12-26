package evaluator

import (
	"explang/parser"
	"fmt"
)

type evaluator struct {
	vars map[string]parser.Node
}

func newEvaluator() *evaluator {
	return &evaluator{
		vars: make(map[string]parser.Node),
	}
}

func (e *evaluator) VisitFile(node parser.File) (parser.Node, error) {
	newFile := &parser.File{Nodes: make([]parser.Node, 0)}
	for _, node := range node.Nodes {
		newNode, err := node.Accept(e)
		if err != nil {
			return nil, err
		}
		newFile.Nodes = append(newFile.Nodes, newNode)
	}
	return newFile, nil
}

func (e *evaluator) VisitDefinition(node parser.Definition) (parser.Node, error) {
	value, err := node.Expression.Accept(e)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate expression: %w", err)
	}
	e.vars[node.Name] = value
	return node, nil
}

func (e *evaluator) VisitString(node parser.String) (parser.Node, error) {
	return node, nil
}

func (e *evaluator) VisitRef(node parser.Ref) (parser.Node, error) {
	val, ok := e.vars[node.Name]
	if !ok {
		return nil, fmt.Errorf("no such variable: %s", node.Name)
	}
	return val, nil
}

func Evaluate(node parser.Node) (parser.Node, error) {
	evaluator := newEvaluator()
	return node.Accept(evaluator)
}
