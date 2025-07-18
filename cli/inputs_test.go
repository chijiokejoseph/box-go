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
			case errors.Is(err, ErrQuit):
				continue
			default:
				t.Errorf("expected err '%s', for string '%s', in '%v'. Got '%s'", ErrQuit, valueIn, quitWords, err)
			}
		}
	}
}

func TestInputNum(t *testing.T) {
	type data struct {
		expected int
		input    string
	}

	dataset := []data{
		{expected: 25, input: "25"},
		{expected: 0, input: "2ab"},
		{expected: 0, input: "Q"},
		{expected: 30, input: " 30"},
	}

	for _, row := range dataset {
		reader := bufio.NewReader(strings.NewReader(row.input))
		actual, err := InputNum(reader, "Enter a number: ")
		switch {
		case errors.Is(err, ErrQuit):
			continue
		case errors.Is(err, ErrSysInput):
			continue
		case errors.Is(err, ErrParseInt):
			continue
		case err != nil:
			{
				t.Errorf("expected err '%s' or '%s' or '%s', got err '%s\n'", ErrQuit, ErrSysInput, ErrParseInt, err)
			}
		}
		if actual != row.expected {
			t.Errorf("expected '%d', got '%d'\n", row.expected, actual)
		}
	}
}


func TestInputOption(t *testing.T) {
	modes := []string{
		"Beginner",
		"Amateur",
		"Regular",
		"Professional",
		"Expert",
	}

	type option struct {
		idx int
		value string
	}

	compare := func(lhs, rhs option) bool {
		return lhs.idx == rhs.idx && lhs.value == rhs.value
	}

	type data struct {
		expected option
		input string
	}

	dataset := []data{
		{expected: option{1, modes[0]}, input: "1"},
		{expected: option{3, modes[2]}, input: "3"},
		{expected: option{}, input: "quit"},
		{expected: option{}, input: "2ab"},
	}

	for _, row := range dataset {
		reader := bufio.NewReader(strings.NewReader(row.input))
		actualIdx, actualValue, err := InputOption(reader, modes, "Select a game mode: ")
		switch {
		case errors.Is(err, ErrQuit) || errors.Is(err, ErrSysInput) || errors.Is(err, ErrParseInt): continue
		case err != nil: {
			t.Errorf("expected err '%s', '%s', or '%s'. Got err '%s'", ErrQuit, ErrSysInput, ErrParseInt, err)
		}
		}
		actual := option{actualIdx, actualValue}
		if !compare(actual, row.expected) {
			t.Errorf("expected '%d. %s', got '%d. %s'\n", row.expected.idx, row.expected.value, actual.idx, actual.value)
		}
	}
}
