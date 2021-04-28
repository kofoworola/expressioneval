package main

import (
	"flag"
	"fmt"
	"log"
)

var expression string

func init() {
	flag.StringVar(&expression, "expression", "", "the expression to evaluate")
}

func main() {
	flag.Parse()
	if expression == "" {
		log.Fatal("expression can not be empty")
	}
	fmt.Println(expression)
}
