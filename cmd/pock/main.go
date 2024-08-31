package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	pock "github.com/loderunner/pocklang"
)

const version = "0.0.0"

func main() {
	fmt.Printf("Pock v%s\n", version)

	rl, err := readline.New("> ")
	if err != nil {
		panic(err.Error())
	}
	interpreter := pock.NewInterpreter()
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
