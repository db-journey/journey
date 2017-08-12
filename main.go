package main

import (
	"log"
	"os"

	_ "github.com/db-journey/bash-driver"
	_ "github.com/db-journey/cassandra-driver"
	_ "github.com/db-journey/crate-driver"
	journey "github.com/db-journey/journey/commands"
	_ "github.com/db-journey/mysql-driver"
	_ "github.com/db-journey/postgresql-driver"
	_ "github.com/db-journey/sqlite3-driver"
	_ "github.com/nexdrew/cockroachdb-driver"
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
	app.Version = "2.0.0"

	app.Flags = journey.Flags()

	app.Commands = journey.Commands()
	return app
}
