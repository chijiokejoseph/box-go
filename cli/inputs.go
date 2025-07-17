package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

var QuitErr = errors.New("player has quit operation")
var SysInputErr = errors.New("system input operation failed")
var ParseIntErr = errors.New("failed to parse str input to int")
var promptWords = [3]string{
	"enter",
	"choose",
	"select",
}
var quitWords = [4]string{
	"quit",
	"q",
	"exit",
	"e",
}

const quitPrompt = "\nEnter 'q', 'quit', 'e', or 'exit' to stop the current operation.\n"

// modPrompt modifies an input prompt string.
// It does so by adding a clause to the string.
// This informs the user on how to exit the
// current operation
//
//   - prompt - the string to be printed out
//     to the user
//     with the aim of requesting user input
//
// Returns: string with the prompt clause for
// cancelling the current operation added
func modPrompt(prompt string) string {
	promptLower := strings.ToLower(prompt) // get `prompt` in lowercase
	for _, word := range promptWords {
		// find the location of the final/main user instruction
		idx := strings.LastIndex(promptLower, word)
		if idx == -1 {
			continue
		}
		// split the prompt at the start of the final/main user instruction
		// and add the quit clause there.
		pre, post := prompt[:idx], prompt[idx:]
		prompt = fmt.Sprintf("%s%s%s", pre, quitPrompt, post)
		break
	}
	return prompt
}

func isTrimmable(str string) bool {
	if len(str) == 0 {
		return false
	}
	return str[0] == ' ' || str[len(str)-1] == ' '
}

// Input prints a prompt to the user, then reads user input
// through the reader argument.
//
//   - reader: a bufio reader created from os.Stdin
//   - prompt: the message printed to the user
//
// Returns:
//
// (string) - the read string with trailing whitespaces and newlines
// trimmed.
//
// (error) - The error that occurs in the program. Note that an end
// of file error (io.EOF) is not considered an error in this
// function
func Input(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Printf("%s: ", modPrompt(prompt))
	resp, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("%w. %w", SysInputErr, err)
	}
	resp = strings.Trim(resp, "\n")
	for isTrimmable(resp) {
		resp = strings.Trim(resp, " ")
	}
	if slices.Contains(quitWords[:], strings.ToLower(resp)) {
		return "", QuitErr
	}
	return resp, nil
}

func InputNum(reader *bufio.Reader, prompt string) (int, error) {
	resp, err := Input(reader, prompt)
	if err != nil {
		return 0, err
	}
	var value int
	_, err = fmt.Sscanf(resp, "%d", &value)
	if err != nil {
		return 0, fmt.Errorf("%w. %w", ParseIntErr, err)
	}
	return value, nil
}

func InputOption(reader *bufio.Reader, options []string, prompt string, title ...string) (idx int, option string, err error) {
	menuTitle := "Select from the following: "
	if len(title) == 1 {
		menuTitle = title[0]
	}
	fmt.Printf("%s\n", menuTitle)
	for i := range options {
		fmt.Printf("%d. %s\n", i+1, options[i])
	}
	fmt.Println()
	resp, err := InputNum(reader, prompt)
	if err != nil {
		return 0, "", err
	}
	idx = resp - 1
	option = options[idx]
	return idx, option, nil
}
