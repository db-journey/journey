package journey

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/db-journey/migrate"
	"github.com/urfave/cli"
)

var MigrateFlags = []cli.Flag{}

//Commands returns the application cli commands:
//create <name>  Create a new migration
//up             Apply all -up- migrations
//down           Apply all -down- migrations
//reset          Down followed by Up
//redo           Roll back most recent migration, then apply it again
//version        Show current migration version
//migrate <n>    Apply migrations -n|+n
//goto <v>       Migrate to version v
func MigrateCommands() cli.Commands {
	return cli.Commands{
		createCommand,
		upCommand,
		downCommand,
		resetCommand,
		redoCommand,
		versionCommand,
		migrateCommand,
		gotoCommand,
	}
}

var createCommand = cli.Command{
	Name:      "create",
	Aliases:   []string{"c"},
	Usage:     "Create a new migration",
	ArgsUsage: "<name>",
	Flags:     MigrateFlags,
	Action: func(ctx *cli.Context) error {
		name := ctx.Args().First()
		if name == "" {
			log.Fatal("Please specify a name for the new migration")
		}
		// if more than one param is passed, create a concat name
		if ctx.NArg() != 1 {
			name = strings.Join(ctx.Args(), "_")
		}

		migrationFile, err := migrate.Create(ctx.GlobalString("url"), ctx.GlobalString("path"), name)
		if err != nil {
			logErr(err).Fatal("Migration failed")
		}

		log.WithFields(log.Fields{
			"up":   migrationFile.UpFile.FileName,
			"down": migrationFile.DownFile.FileName,
		}).Infof("Version %v migration files created in %v:\n", migrationFile.Version, ctx.GlobalString("path"))
		return nil
	},
}

var upCommand = cli.Command{
	Name:  "up",
	Usage: "Apply all -up- migrations",
	Flags: MigrateFlags,
	Action: func(ctx *cli.Context) error {
		log.Info("Applying all -up- migrations")
		err := migrate.Up(ctx.GlobalString("url"), ctx.GlobalString("path"))
		if err != nil {
			logErr(err).Fatal("Failed to apply all -up- migrations")
		}
		logCurrentVersion(ctx.GlobalString("url"), ctx.GlobalString("path"))
		return nil
	},
}

var downCommand = cli.Command{
	Name:  "down",
	Usage: "Apply all -down- migrations",
	Flags: MigrateFlags,
	Action: func(ctx *cli.Context) error {
		log.Info("Applying all -down- migrations")
		err := migrate.Down(ctx.GlobalString("url"), ctx.GlobalString("path"))
		if err != nil {
			logErr(err).Fatal("Failed to apply all -down- migrations")
		}
		logCurrentVersion(ctx.GlobalString("url"), ctx.GlobalString("path"))
		return nil
	},
}

var redoCommand = cli.Command{
	Name:    "redo",
	Aliases: []string{"r"},
	Usage:   "Roll back most recent migration, then apply it again",
	Flags:   MigrateFlags,
	Action: func(ctx *cli.Context) error {
		log.Info("Redoing last migration")
		err := migrate.Redo(ctx.GlobalString("url"), ctx.GlobalString("path"))
		logErr(err).Fatal("Failed to redo last migration")
		logCurrentVersion(ctx.GlobalString("url"), ctx.GlobalString("path"))
		return nil
	},
}

var versionCommand = cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "Show current migration version",
	Flags:   MigrateFlags,
	Action: func(ctx *cli.Context) error {
		version, err := migrate.Version(ctx.GlobalString("url"), ctx.GlobalString("path"))
		if err != nil {
			logErr(err).Fatal("Unable to fetch version")
		}

		log.Infof("Current version: %d", version)
		return nil
	},
}

var resetCommand = cli.Command{
	Name:  "reset",
	Usage: "Down followed by Up",
	Flags: MigrateFlags,
	Action: func(ctx *cli.Context) error {
		log.Info("Reseting database")
		err := migrate.Redo(ctx.GlobalString("url"), ctx.GlobalString("path"))
		if err != nil {
			logErr(err).Fatal("Failed to reset database")
		}
		logCurrentVersion(ctx.GlobalString("url"), ctx.GlobalString("path"))
		return nil
	},
}

var migrateCommand = cli.Command{
	Name:            "migrate",
	Aliases:         []string{"m"},
	Usage:           "Apply migrations -n|+n",
	ArgsUsage:       "<n>",
	Flags:           MigrateFlags,
	SkipFlagParsing: true,
	Action: func(ctx *cli.Context) error {
		relativeN := ctx.Args().First()
		relativeNInt, err := strconv.Atoi(relativeN)
		if err != nil {
			logErr(err).Fatal("Unable to parse param <n>")
		}

		log.Infof("Applying %d migrations", relativeNInt)

		err = migrate.Migrate(ctx.GlobalString("url"), ctx.GlobalString("path"), relativeNInt)
		if err != nil {
			logErr(err).Fatalf("Failed to apply %d migrations", relativeNInt)
		}
		logCurrentVersion(ctx.GlobalString("url"), ctx.GlobalString("path"))
		return nil
	},
}

var gotoCommand = cli.Command{
	Name:      "goto",
	Aliases:   []string{"g"},
	Usage:     "Migrate to version <v>",
	ArgsUsage: "<v>",
	Flags:     MigrateFlags,
	Action: func(ctx *cli.Context) error {
		toVersion := ctx.Args().First()
		toVersionInt, err := strconv.Atoi(toVersion)
		if err != nil || toVersionInt < 0 {
			logErr(err).Fatal("Unable to parse param <v>")
		}

		log.Infof("Migrating to version %d", toVersionInt)

		currentVersion, err := migrate.Version(ctx.GlobalString("url"), ctx.GlobalString("path"))
		if err != nil {
			logErr(err).Fatalf("failed to migrate to version %d", toVersionInt)
		}

		relativeNInt := toVersionInt - int(currentVersion)

		err = migrate.Migrate(ctx.GlobalString("url"), ctx.GlobalString("path"), relativeNInt)
		if err != nil {
			logErr(err).Fatalf("Failed to migrate to vefrsion %d", toVersionInt)
		}
		logCurrentVersion(ctx.GlobalString("url"), ctx.GlobalString("path"))
		return nil
	},
}

func logErr(err error) *log.Entry {
	return log.WithError(err)
}

func logCurrentVersion(url, migrationsPath string) {
	version, err := migrate.Version(url, migrationsPath)
	if err != nil {
		logErr(err).Fatal("Unable to fetch version")
	}
	log.WithField("current-version", version).Info("done")
}
