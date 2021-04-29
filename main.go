package main

import (
	"errors"
	"fmt"
	"os"
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

var (
	symbolRegex = `[+\*\/\-\(\)]`

	reg *regexp.Regexp
)

func init() {
	reg = regexp.MustCompile(symbolRegex)
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

func main() {
	args := os.Args[1:]
	if len(args) < 1 || args[0] == "" {
		panic("missing expression")
	}

	expression := args[0]
	sol, err := evaluate(expression)
	if err != nil {
		panic(err)
	}
	fmt.Println(sol)
}

func evaluate(expression string) (float64, error) {
	tokens := generateTokens(expression)

	rpn, err := convertToRPN(tokens)
	if err != nil {
		return 0, err
	}

	// evaluate the RPN
	st := stack{[]string{}}
	for _, tok := range rpn {
		if !reg.MatchString(tok) {
			st.push(tok)
			continue
		}
		first, err := strconv.ParseFloat(st.pop(), 64)
		if err != nil {
			return 0, err
		}
		second, err := strconv.ParseFloat(st.pop(), 64)
		if err != nil {
			return 0, err
		}

		var eval float64
		switch tok {
		case "+":
			eval = second + first
		case "-":
			eval = second - first
		case "*":
			eval = second * first
		case "/":
			eval = second / first
		default:
			return 0, errors.New("unknown operator")
		}
		st.push(fmt.Sprintf("%f", eval))
	}

	return strconv.ParseFloat(st.pop(), 64)
}

// generateTokens takes and expression and splits into tokens
func generateTokens(expr string) []string {
	var tokens []string
	var lastDigits strings.Builder
	for _, c := range expr {
		if c == ' ' {
			continue
		}
		// check if is a digit and add to the digit builder
		if unicode.IsDigit(c) {
			lastDigits.WriteRune(c)
		} else {
			// if not a digit then empty out the digits into the token list
			// and add the operator
			lastDigit := lastDigits.String()
			if lastDigit != "" {
				tokens = append(tokens, lastDigit)
			}
			lastDigits.Reset()
			tokens = append(tokens, string(c))
		}
	}
	// empty the digit builder after last operation
	if str := lastDigits.String(); str != "" {
		tokens = append(tokens, str)
	}
	return tokens
}

// convertToRPN converts the token array into the reverse polish notation
// see https://en.wikipedia.org/wiki/Reverse_Polish_notation
func convertToRPN(tokens []string) ([]string, error) {
	var rpn []string
	operators := stack{[]string{}}
	for _, token := range tokens {
		if _, err := strconv.Atoi(token); err == nil {
			rpn = append(rpn, token)
		} else {
			if !reg.MatchString(token) {
				return nil, errors.New("invalid expression")
			}
			if operators.peek() == "" {
				operators.push(token)
				continue
			}
			// handle parenthesis
			if token == "(" {
				operators.push(token)
				continue
			}
			if token == ")" {
				prev := operators.pop()
				for {
					if prev == "(" || prev == "" {
						break
					}
					rpn = append(rpn, prev)
					prev = operators.pop()
				}
				if operators.peek() == "(" {
					operators.pop()
				}
				continue
			}

			// check if previous operator takes precedence
			prev := operators.peek()
			for precedenceMap[prev] > precedenceMap[token] && prev != "" {
				rpn = append(rpn, operators.pop())
				prev = operators.peek()
			}
			operators.push(token)
		}
	}

	// flush all operators into the rpn stack
	for tok := operators.peek(); tok != ""; tok = operators.peek() {
		rpn = append(rpn, operators.pop())
	}
	return rpn, nil
}
