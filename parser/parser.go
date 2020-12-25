package parser

import (
	"errors"
	"explang/tokenizer"
)

type Node interface {
}

type file struct {
	nodes []Node
}

type definition struct {
	name       string
	expression Node
}

type call struct {
}

type ref struct {
	name string
}

type str struct {
	value string
}

type parser struct {
	tokens []tokenizer.Token
}

func (p *parser) Done() bool {
	return len(p.tokens) == 0
}

func (p *parser) Next() tokenizer.Token {
	next := p.tokens[0]
	p.tokens = p.tokens[1:]
	return next
}

func (p *parser) ParseFile() (Node, error) {
	res := &file{nodes: make([]Node, 0)}
	for {
		if p.Done() {
			break
		}
		node, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		res.nodes = append(res.nodes, node)
	}
	return res, nil
}

func (p *parser) ParseExpression() (Node, error) {
	if p.Done() {
		return nil, errors.New("expected expression")
	}
	next := p.Next()
	switch next.Type {
	case tokenizer.LeftParen:
		if p.Done() {
			return nil, errors.New("expected definition or function call")
		}
		next = p.Next()
		if next.Type != tokenizer.Identifier {
			return nil, errors.New("expected 'define' or function name")
		}
		if next.Value == "define" {
			node, err := p.ParseDefinition()
			if err != nil {
				return nil, errors.New("could not parse definition")
			}
			return node, nil
		} else {
			return nil, errors.New("unsupported function call")
		}

	case tokenizer.String:
		return str{value: next.Value}, nil

	case tokenizer.Identifier:
		return ref{name: next.Value}, nil

	default:
		return nil, errors.New("unsupported expression")
	}

}

func (p *parser) ParseDefinition() (Node, error) {
	if p.Done() {
		return nil, errors.New("expected definition name")
	}
	next := p.Next()
	if next.Type != tokenizer.Identifier {
		return nil, errors.New("expected identifier for definition name")
	}
	name := next.Value
	node, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if p.Done() || p.Next().Type != tokenizer.RightParen {
		return nil, errors.New("expected right paren")
	}

	return definition{name: name, expression: node}, nil
}

func Parse(tokens []tokenizer.Token) (Node, error) {
	parser := &parser{tokens: tokens}
	return parser.ParseFile()
}
