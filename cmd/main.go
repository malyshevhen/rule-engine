package main

import (
	"fmt"
	"os"

	"github.com/malyshevhen/rule-engine/cmd/app"
)

func run() error {
	app := app.New()
	return app.Run()
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
