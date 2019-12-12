package main

import (
	"log"
	"os"

	_ "github.com/db-journey/bash-driver"
	_ "github.com/db-journey/cassandra-driver"
	_ "github.com/db-journey/crate-driver"
	journey "github.com/db-journey/journey/v2/commands"
	_ "github.com/db-journey/mysql-driver"
	_ "github.com/db-journey/postgresql-driver"
	_ "github.com/db-journey/sqlite3-driver"
	"github.com/urfave/cli"
)

func main() {
	app := App()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func App() *cli.App {
	app := cli.NewApp()
	app.Usage = "Migrations and cronjobs for databases"
	app.Version = "2.1.1"

	app.Flags = journey.Flags()

	app.Commands = journey.Commands()
	return app
}
