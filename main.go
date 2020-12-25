package main

import (
	"explang/parser"
	"explang/tokenizer"
	"fmt"
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
	fmt.Printf("tokens: %+v\n", tokens)
	res, err := parser.Parse(tokens)
	if err != nil {
		log.Fatalf("parsing failed: %v", err)
	}
	fmt.Printf("nodes: %+v", res)
}
