package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "   hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "    heLLo     WOrld    !    ",
			expected: []string{"hello", "world", "!"},
		},
		{
			input:    "HELLOWORLD",
			expected: []string{"helloworld"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}
	fmt.Println("TestCleanInput:...")
	passed := 0
	failed := 0
	for _, c := range cases {
		actual := cleanInput(c.input)
		fmt.Printf("actual: %#v, expected: %#v\n", actual, c.expected)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			failed++
			t.Errorf("length of return value for '%s':%#v != %#v", c.input, actual, c.expected)
		} else {
			passed++
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				failed++
				t.Errorf("return value '%s' != '%s'", word, expectedWord)
			} else {
				passed++
			}
		}
	}
	fmt.Println("+++++++++++++++++++++++++++++++++")
	fmt.Printf("passed: %d | failed: %d\n", passed, failed)
	fmt.Println("+++++++++++++++++++++++++++++++++")
}
