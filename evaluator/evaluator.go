package evaluator

import (
	"cacophony/parser"
	"cacophony/tokenizer"
	"errors"
	"fmt"
	"io"
	"log"
)

type evaluator struct {
	*environment
	writer io.Writer
}

func newEvaluator(writer io.Writer) *evaluator {
	return &evaluator{
		environment: newEnvironment(),
		writer:      writer,
	}
}

func (e *evaluator) VisitFile(node parser.File) (parser.Node, error) {
	newFile := &parser.File{Nodes: make([]parser.Node, 0)}
	for _, node := range node.Nodes {
		newNode, err := node.Accept(e)
		if err != nil {
			return nil, err
		}
		if !newNode.IsReducible() {
			fmt.Fprintf(e.writer, "%s\n", newNode)
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
	e.Set(node.Name, value)
	return node, nil
}

func (e *evaluator) VisitIf(node parser.If) (parser.Node, error) {
	cond, err := node.Cond.Accept(e)
	if err != nil {
		return nil, err
	}
	condAsBool, ok := cond.(parser.Boolean)
	if !ok {
		return nil, errors.New("expected boolean expression")
	}
	if condAsBool.Value {
		return node.TrueBranch.Accept(e)
	} else {
		return node.FalseBranch.Accept(e)
	}
}

func (e *evaluator) VisitString(node parser.String) (parser.Node, error) {
	return node, nil
}

func (e *evaluator) VisitBoolean(node parser.Boolean) (parser.Node, error) {
	return node, nil
}

func (e *evaluator) VisitRef(node parser.Ref) (parser.Node, error) {
	val, ok := e.Get(node.Name)
	if !ok {
		return nil, fmt.Errorf("no such variable: %s", node.Name)
	}
	return val, nil
}

func evaluate(node parser.Node, writer io.Writer) (parser.Node, error) {
	evaluator := newEvaluator(writer)
	return node.Accept(evaluator)
}

func Evaluate(contents string, writer io.Writer) (parser.Node, error) {
	tokens := tokenizer.Tokenize(contents)
	log.Printf("tokens: %+v", tokens)
	res, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}
	log.Printf("nodes: %+v", res)

	node, err := evaluate(res, writer)
	if err != nil {
		return nil, err
	}
	return node, nil
}
