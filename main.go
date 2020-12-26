package main

import (
	"explang/evaluator"
	"explang/parser"
	"explang/tokenizer"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("you must provide a file")
	}

	reader, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not read from file: %v", err)
	}
	tokens, err := tokenizer.Tokenize(reader)
	if err != nil {
		log.Fatalf("tokenization failed: %v", err)
	}
	log.Printf("tokens: %+v", tokens)
	res, err := parser.Parse(tokens)
	if err != nil {
		log.Fatalf("parsing failed: %v", err)
	}
	log.Printf("nodes: %+v", res)

	node, err := evaluator.Evaluate(res)
	if err != nil {
		log.Fatalf("evaluation failed: %v", err)
	}
	log.Printf("evaluated nodes: %+v", node)
}
