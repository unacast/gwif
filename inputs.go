package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Ask(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n] ", question)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		}
		if response == "n" || response == "no" {
			return false
		}
		fmt.Println("Please answer yes or no.")
	}
}

func RequiredAsk(question, errorMessage string) bool {
	for {
		if Ask(question) {
			return true
		}
		fmt.Println(errorMessage)
	}
}

func GetInput(prompt string) string {
	fmt.Print("\n", prompt, " ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return ""
		}
		return ""
	}
	return scanner.Text()
}
