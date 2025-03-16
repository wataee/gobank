package main

import (
	"gobank1/pkg/database"
	"gobank1/pkg/cli"

)

func main() {

	pool, _ := database.ConnectDB("")
	cli.RunCLI()

	defer database.CloseDB(pool)

}