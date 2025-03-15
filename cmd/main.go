package main

import (
	"gobank1/pkg/database"
	"gobank1/pkg/cli"

)

func main() {

	pool, _ := database.ConnectDB("postgres://postgres:1234@localhost:5432/postgres")
	cli.RunCLI()

	defer database.CloseDB(pool)

}