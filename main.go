package main

import (
	"cacophony/evaluator"
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
	node, err := evaluator.Evaluate(reader, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("evaluated nodes: %+v", node)
}
