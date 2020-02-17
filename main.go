package main

import (
	"log"
	"os"

	journey "github.com/db-journey/journey/v2/commands"
	_ "github.com/db-journey/migrate/v2/drivers/bash-driver"
	_ "github.com/db-journey/migrate/v2/drivers/cassandra-driver"
	_ "github.com/db-journey/migrate/v2/drivers/crate-driver"
	_ "github.com/db-journey/migrate/v2/drivers/mysql-driver"
	_ "github.com/db-journey/migrate/v2/drivers/postgresql-driver"
	_ "github.com/db-journey/migrate/v2/drivers/sqlite3-driver"
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
	app.Version = "2.2.5"

	app.Flags = journey.Flags()

	app.Commands = journey.Commands()
	return app
}
