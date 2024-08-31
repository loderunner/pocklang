package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	pock "github.com/loderunner/pocklang"
)

const version = "0.0.0"

var statePath = flag.String("state", "", "a JSON file to be loaded as interpreter state")

func main() {
	flag.Parse()
	var interpreter *pock.Interpreter
	if *statePath == "" {
		interpreter = pock.NewInterpreter()
	} else {
		var state map[string]any
		buf, err := os.ReadFile(*statePath)
		if err != nil {
			panic(err.Error())
		}
		err = json.Unmarshal(buf, &state)
		if err != nil {
			panic(err.Error())
		}
		interpreter, err = pock.NewInterpreterWithState(state)
		if err != nil {
			panic(err.Error())
		}
	}

	fmt.Printf("Pock v%s\n", version)

	rl, err := readline.New("> ")
	if err != nil {
		panic(err.Error())
	}
	for {
		src, err := rl.Readline()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			panic(err.Error())
		}

		tokens, err := pock.Scan(strings.NewReader(src))
		if err != nil {
			fmt.Printf("scan error: %s\n", err)
			continue
		}

		expr, err := pock.Parse(tokens)
		if err != nil {
			fmt.Printf("parse error: %s\n", err)
			continue
		}

		value, err := interpreter.Evaluate(expr)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}

		switch value.(type) {
		case string:
			fmt.Printf("\"%s\"", value)
		case nil:
			fmt.Print("null")
		default:
			fmt.Print(value)
		}
		fmt.Println()
	}
}
