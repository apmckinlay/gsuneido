// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/apmckinlay/gsuneido/llm"
)

const model = "openai/o4-mini"

var done chan struct{}

func main() {
	agent := llm.NewAgent("https://openrouter.ai/api/v1",
		os.Getenv("OPENROUTER_API_KEY"), model, "", outfn)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Chat with", model, "('q' to quit)")
	fmt.Println()

	for {
		fmt.Print("You: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("read input: ", err)
		}
		fmt.Println()

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if input == "q" {
			break
		}

		fmt.Print("LLM: ")

		done = make(chan struct{})
		agent.Input(input)
		<-done

		fmt.Println()
		fmt.Println()
	}
}

// outfn handles streaming output from the agent.
func outfn(what, data string, _ *llm.ToolApproval) {
	switch what {
	case "think":
		fmt.Printf("\x1b[34m%s\x1b[0m", data) // blue
	case "output":
		fmt.Print(data)
	case "tool":
		fmt.Printf("\x1b[32m%s\x1b[0m", data) // green
	case "complete":
		close(done)
	}
}
