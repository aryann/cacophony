package main

import (
	"cacophony/evaluator"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("you must provide a file")
	}

	contents, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("could not read from file: %v", err)
	}
	node, err := evaluator.Evaluate(string(contents), os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	log.Printf("evaluated nodes: %+v", node)
}
