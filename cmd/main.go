// Package main Rule Engine API
//
//	@title			Rule Engine API
//	@version		1.0
//	@description	A robust rule engine microservice for IoT automation. Allows users to create and manage custom automation rules with Lua script execution in a secure sandboxed environment.
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host		localhost:8080
//	@BasePath	/api/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				API key authentication. Format: "ApiKey <your-api-key>"
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT Bearer token authentication. Format: "Bearer <your-jwt-token>"
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
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
