package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

var ErrQuit = errors.New("player has quit operation")
var ErrSysInput = errors.New("system input operation failed")
var ErrParseInt = errors.New("failed to parse str input to int")
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
//   - prompt: the string to be printed out
//     to the user
//     with the aim of requesting user input
//
// Returns: 
// 	string with the prompt clause for
// 	cancelling the current operation added
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
// 	(string) - the read string with trailing whitespaces and newlines
// 		trimmed.
// 	(error) - The error that occurs in the program. Note that an end
// 		of file error (io.EOF) is not considered an error in this
// 		function
func Input(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Printf("%s: ", modPrompt(prompt))
	resp, err := reader.ReadString('\n')

	// check for non-nil errors that are not io.EOF
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("%w. %w", ErrSysInput, err)
	}

	// trim all whitespaces in the user's response
	resp = strings.Trim(resp, "\n")
	for isTrimmable(resp) {
		resp = strings.Trim(resp, " ")
	}

	// handle quit or exit messages
	if slices.Contains(quitWords[:], strings.ToLower(resp)) {
		return "", ErrQuit
	}

	return resp, nil
}

// InputNum prompts a user to enter a number.
// It returns the number and the error encountered
// when parsing the user input to a valid int.
//
// 	- reader: a bufio reader that reads user input
// 	- prompt: the instruction printed out to the user
// 
// Returns:
// 
// 	(int) - the number entered from the console by the user
// 	(error) - the error encountered when parsing user input
func InputNum(reader *bufio.Reader, prompt string) (int, error) {
	// get user input
	resp, err := Input(reader, prompt)
	if err != nil {
		return 0, err
	}
	
	// convert user input to int
	value, err := strconv.Atoi(resp)
	if err != nil {
		return 0, fmt.Errorf("%w. %w", ErrParseInt, err)
	}

	return value, nil
}


// InputOption prompts the user to select an option
// from a list of available choices.
//
// 	- reader: a reader that reads user input from stdin
// 	- options: holds the list of prompt
// 	- prompt: the instructions to be displayed to the user
// 	- title: a semi-optional field that can be empty. when
// 		given a single argument, it specifies the `menuTitle`
// 		variable in the function.
//
// Returns:
//
// 	(num int): the position of the selected option starting
//		from 1
// 	(option string): the corresponding option selected by 
//		the user based off his input
// 	(err error): the error encountered while executing the
// 		the function
func InputOption(reader *bufio.Reader, options []string, prompt string, title ...string) (num int, option string, err error) {
	// set menuTitle to default and update if specified
	// in the function call
	menuTitle := "Select from the following: "
	if len(title) == 1 {
		menuTitle = title[0]
	}
	
	fmt.Printf("%s\n", menuTitle)

	// display options
	for i := range options {
		fmt.Printf("%d. %s\n", i+1, options[i])
	}
	fmt.Println()

	// read user option choice as int
	num, err = InputNum(reader, prompt)
	if err != nil {
		return 0, "", err
	}

	option = options[num - 1]
	return num, option, nil
}
