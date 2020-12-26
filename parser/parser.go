package parser

import (
	"cacophony/tokenizer"
	"errors"
	"fmt"
	"strconv"
)

type Visitor interface {
	VisitFile(node File) (Node, error)
	VisitDefinition(node Definition) (Node, error)
	VisitIf(node If) (Node, error)
	VisitString(node String) (Node, error)
	VisitBoolean(node Boolean) (Node, error)
	VisitRef(node Ref) (Node, error)
}

type Node interface {
	Accept(visitor Visitor) (Node, error)
	IsReducible() bool
}

type File struct {
	Nodes []Node
}

func (f File) Accept(visitor Visitor) (Node, error) { return visitor.VisitFile(f) }
func (File) IsReducible() bool                      { return true }

type Definition struct {
	Name       string
	Expression Node
}

func (d Definition) Accept(visitor Visitor) (Node, error) { return visitor.VisitDefinition(d) }
func (Definition) IsReducible() bool                      { return true }

type If struct {
	Cond        Node
	TrueBranch  Node
	FalseBranch Node
}

func (i If) Accept(visitor Visitor) (Node, error) { return visitor.VisitIf(i) }
func (If) IsReducible() bool                      { return true }

type Ref struct {
	Name string
}

func (r Ref) Accept(visitor Visitor) (Node, error) { return visitor.VisitRef(r) }
func (Ref) IsReducible() bool                      { return true }

type String struct {
	Value string
}

func (s String) Accept(visitor Visitor) (Node, error) { return visitor.VisitString(s) }
func (String) IsReducible() bool                      { return false }
func (s String) String() string                       { return strconv.Quote(s.Value) }

type Boolean struct {
	Value bool
}

func (b Boolean) Accept(visitor Visitor) (Node, error) { return visitor.VisitBoolean(b) }
func (Boolean) IsReducible() bool                      { return false }
func (b Boolean) String() string                       { return ":" + strconv.FormatBool(b.Value) }

type parser struct {
	tokens []tokenizer.Token
}

func (p parser) Done() bool {
	return len(p.tokens) == 0
}

func (p *parser) Next() tokenizer.Token {
	next := p.tokens[0]
	p.tokens = p.tokens[1:]
	return next
}

func (p *parser) ParseFile() (Node, error) {
	res := &File{Nodes: make([]Node, 0)}
	for {
		if p.Done() {
			break
		}
		node, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		res.Nodes = append(res.Nodes, node)
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

		switch next.Type {
		case tokenizer.BuiltIn:
			switch next.Value {
			case "define":
				node, err := p.ParseDefinition()
				if err != nil {
					return nil, errors.New("could not parse definition")
				}
				return node, nil

			case "if":
				node, err := p.ParseIf()
				if err != nil {
					return nil, errors.New("could not parse if")
				}
				return node, nil

			case "true", "false":
				return nil, fmt.Errorf("'%s' cannot be called", next.Value)
			default:
				return nil, fmt.Errorf("unknown keyword '%s'", next.Value)
			}

		case tokenizer.Identifier:
			return nil, errors.New("function call are not yet supported")
		default:
			return nil, errors.New("expected keyword or function name")
		}

	case tokenizer.String:
		return String{Value: next.Value}, nil

	case tokenizer.Identifier:
		return Ref{Name: next.Value}, nil

	case tokenizer.BuiltIn:
		switch next.Value {
		case "true":
			return Boolean{Value: true}, nil
		case "false":
			return Boolean{Value: false}, nil
		default:
			return nil, errors.New("unexpected built-in reference")
		}

	default:
		return nil, errors.New("unsupported expression")
	}
}

func (p *parser) ParseIf() (Node, error) {
	cond, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	trueBranch, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	falseBranch, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if p.Done() || p.Next().Type != tokenizer.RightParen {
		return nil, errors.New("expected right paren")
	}
	return If{
		Cond:        cond,
		TrueBranch:  trueBranch,
		FalseBranch: falseBranch,
	}, nil
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
	return Definition{Name: name, Expression: node}, nil
}

func Parse(tokens []tokenizer.Token) (Node, error) {
	parser := &parser{tokens: tokens}
	return parser.ParseFile()
}
