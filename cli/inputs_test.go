package cli

import (
	"bufio"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestModPrompt(t *testing.T) {
	query := "Enter a number: "
	expected := fmt.Sprintf("%s%s", quitPrompt, query)
	actual := modPrompt(query)
	fmt.Printf("actual = %s\n", actual)
	fmt.Printf("expected = %s\n", expected)
	if actual != expected {
		t.Errorf("\nexpected %s, got %s\n", expected, actual)
	}
}

func TestInput(t *testing.T) {
	type data struct {
		expected string
		input    string
	}
	dataset := []data{
		{expected: "I am hungry", input: "I am hungry "},
		{expected: "", input: " quit "},
		{expected: "", input: " Quit"},
		{expected: "I am sorry", input: "I am sorry "},
		{expected: "", input: "exit "},
		{expected: "", input: "Q "},
		{expected: "I want to go home", input: " I want to go home"},
		{expected: "", input: "E"},
	}

	for _, row := range dataset {
		valueIn, expected := row.input, row.expected
		buf := bufio.NewReader(strings.NewReader(valueIn))
		actual, err := Input(buf, "How are you feeling?: ")
		if actual != expected {
			t.Errorf("expected '%s', got '%s'", expected, actual)
		}
		if slices.Contains(quitWords[:], strings.ToLower(expected)) {
			switch {
			case errors.Is(err, QuitErr):
				continue
			default:
				t.Errorf("expected err '%s', for string '%s', in '%v'. Got '%s'", QuitErr, valueIn, quitWords, err)
			}
		}
	}
}
