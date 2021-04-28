package main

import (
	"flag"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var precedenceMap = map[string]int{
	"/": 4,
	"*": 3,
	"+": 2,
	"-": 1,
}

// we will be using this stack a lot
type stack struct {
	slice []string
}

func (s *stack) push(val string) {
	s.slice = append(s.slice, val)
}

func (s *stack) pop() (val string) {
	if len(s.slice) < 1 {
		return ""
	}
	val = s.slice[len(s.slice)-1]
	s.slice = s.slice[:len(s.slice)-1]
	return
}

func (s *stack) peek() string {
	if len(s.slice) < 1 {
		return ""
	}
	return s.slice[len(s.slice)-1]
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

	convertToRPN(tokens)
}

// support bracket
func convertToRPN(tokens []string) *stack {
	reg := regexp.MustCompile(`[+\*\/\-\(\)]`)
	rpn := stack{[]string{}}
	operators := stack{[]string{}}
	for _, token := range tokens {
		if _, err := strconv.Atoi(token); err == nil {
			rpn.push(token)
		} else {
			if !reg.MatchString(token) {
				panic("invalid expression")
			}
			if operators.peek() == "" {
				operators.push(token)
				continue
			}

			// check if previous operator takes precedence
			prev := operators.peek()
			for precedenceMap[prev] > precedenceMap[token] && prev != "" {
				rpn.push(operators.pop())
				prev = operators.peek()
			}
			operators.push(token)
		}
	}
	for tok := operators.peek(); tok != ""; tok = operators.peek() {
		rpn.push(operators.pop())
	}
	return &rpn
}
