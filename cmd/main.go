package main

import (
	"auth/cmd/app"
	_ "auth/docs"
)

// @title           Auth Service
// @version         0.0.1
// @description service for auth users

// @host      localhost:8080
// @BasePath  /
func main() {
	app.Run()
}
