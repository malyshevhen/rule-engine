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

// @title Rule Engine API
// @version 1.0
// @description This is a rule engine server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
