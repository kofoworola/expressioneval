package main

import (
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {
	s := stack{[]string{}}

	s.push("hi")
	s.push("there")
	if l := len(s.slice); l != 2 {
		t.Fatalf("expected 2 items, got %d", l)
	}

	val := s.peek()
	if val != "there" || len(s.slice) != 2 {
		t.Errorf(
			"expected val 'there' and slice length of 2 got %s and %d respectively",
			val,
			len(s.slice),
		)
	}

	val = s.pop()
	if val != "there" || len(s.slice) != 1 {
		t.Errorf(
			"expected val 'there' and slice length of 1 got %s and %d respectively",
			val,
			len(s.slice),
		)
	}
}

func TestGenerateTokens(t *testing.T) {
	testCases := []struct {
		expression string
		expected   []string
	}{
		{
			expression: "3+4*(12-1)",
			expected:   []string{"3", "+", "4", "*", "(", "12", "-", "1", ")"},
		},
		{
			expression: "3 + 4*2 ",
			expected:   []string{"3", "+", "4", "*", "2"},
		},
	}

	for _, test := range testCases {
		tokens := generateTokens(test.expression)
		if !reflect.DeepEqual(tokens, test.expected) {
			t.Errorf("expected %v but got %v", test.expected, tokens)
		}
	}
}

func TestConvertToRPN(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		expected   []string
	}{
		{
			name:       "NoBrackets",
			expression: "5/2*3",
			expected:   []string{"5", "2", "/", "3", "*"},
		},
		{
			name:       "HasBrackets",
			expression: "(5+2)*10+3",
			expected:   []string{"5", "2", "+", "10", "*", "3", "+"},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tokens := generateTokens(test.expression)
			rpn, err := convertToRPN(tokens)
			if err != nil {
				t.Errorf("error converting to rpn: %v", err)
			}
			if !reflect.DeepEqual(test.expected, rpn) {
				t.Errorf("expected %v, got %v", test.expected, rpn)
			}
		})
	}
}
