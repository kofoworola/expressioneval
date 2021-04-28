package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"unicode"
)

// we will be using this stack a lot
type stack []string

func (s stack) push(val string) {
	s = append(s, val)
}

func (s stack) pop() (val string) {
	if len(s) < 1 {
		return ""
	}
	val = s[len(s)-1]
	s = s[:len(s)-1]
	return
}

func (s stack) peek() string {
	if len(s) < 1 {
		return ""
	}
	return s[len(s)-1]
}

var expression string

func init() {
	flag.StringVar(&expression, "expression", "", "the expression to evaluate")
}

func main() {
	flag.Parse()
	if expression == "" {
		log.Fatal("expression can not be empty")
	}

	// convert to tokens
	var tokens []string
	var lastDigits strings.Builder
	for _, c := range expression {
		if unicode.IsDigit(c) {
			lastDigits.WriteRune(c)
		} else {
			lastDigit := lastDigits.String()
			if lastDigit != "" {
				tokens = append(tokens, lastDigit)
			}
			lastDigits.Reset()
			tokens = append(tokens, string(c))
		}
	}
	// empty the string builder
	if str := lastDigits.String(); str != "" {
		tokens = append(tokens, str)
	}
	fmt.Println(tokens)
}
