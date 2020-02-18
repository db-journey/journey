package commands

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/db-journey/migrate/v2"
	"github.com/db-journey/migrate/v2/file"
	"github.com/urfave/cli/v2"
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
		applyCommand,
		rollbackCommand,
		gotoCommand,
	}
}

var createCommand = &cli.Command{
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
			name = strings.Join(ctx.Args().Slice(), "_")
		}

		migrate, _, cancel := newMigrateWithCtx(ctx.String("url"), ctx.String("path"))
		defer cancel()

		migrationFile, err := migrate.Create(name)
		if err != nil {
			logErr(err).Fatal("Migration failed")
		}

		log.WithFields(log.Fields{
			"up":   migrationFile.UpFile.FileName,
			"down": migrationFile.DownFile.FileName,
		}).Infof("Version %v migration files created in %v:\n", migrationFile.Version, ctx.String("path"))
		return nil
	},
}

var upCommand = &cli.Command{
	Name:  "up",
	Usage: "Apply all -up- migrations",
	Flags: MigrateFlags,
	Action: func(ctx *cli.Context) error {
		log.Info("Applying all -up- migrations")
		migrate, mctx, cancel := newMigrateWithCtx(ctx.String("url"), ctx.String("path"))
		defer cancel()
		err := migrate.Up(mctx)
		if err != nil {
			logErr(err).Fatal("Failed to apply all -up- migrations")
		}
		logCurrentVersion(mctx, migrate)
		return nil
	},
}

var downCommand = &cli.Command{
	Name:  "down",
	Usage: "Apply all -down- migrations",
	Flags: MigrateFlags,
	Action: func(ctx *cli.Context) error {
		log.Info("Applying all -down- migrations")
		migrate, mctx, cancel := newMigrateWithCtx(ctx.String("url"), ctx.String("path"))
		defer cancel()
		err := migrate.Down(mctx)
		if err != nil {
			logErr(err).Fatal("Failed to apply all -down- migrations")
		}
		logCurrentVersion(mctx, migrate)
		return nil
	},
}

var redoCommand = &cli.Command{
	Name:    "redo",
	Aliases: []string{"r"},
	Usage:   "Roll back most recent migration, then apply it again",
	Flags:   MigrateFlags,
	Action: func(ctx *cli.Context) error {
		log.Info("Redoing last migration")
		migrate, mctx, cancel := newMigrateWithCtx(ctx.String("url"), ctx.String("path"))
		defer cancel()
		err := migrate.Redo(mctx)
		if err != nil {
			logErr(err).Fatal("Failed to redo last migration")
		}
		logCurrentVersion(mctx, migrate)
		return nil
	},
}

var versionCommand = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "Show current migration version",
	Flags:   MigrateFlags,
	Action: func(ctx *cli.Context) error {
		migrate, mctx, cancel := newMigrateWithCtx(ctx.String("url"), ctx.String("path"))
		defer cancel()
		version, err := migrate.Version(mctx)
		if err != nil {
			logErr(err).Fatal("Unable to fetch version")
		}

		log.Infof("Current version: %d", version)
		return nil
	},
}

var resetCommand = &cli.Command{
	Name:  "reset",
	Usage: "Down followed by Up",
	Flags: MigrateFlags,
	Action: func(ctx *cli.Context) error {
		log.Info("Reseting database")
		migrate, mctx, cancel := newMigrateWithCtx(ctx.String("url"), ctx.String("path"))
		defer cancel()
		err := migrate.Redo(mctx)
		if err != nil {
			logErr(err).Fatal("Failed to reset database")
		}
		logCurrentVersion(mctx, migrate)
		return nil
	},
}

var migrateCommand = &cli.Command{
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

		migrate, mctx, cancel := newMigrateWithCtx(ctx.String("url"), ctx.String("path"))
		defer cancel()
		err = migrate.Migrate(mctx, relativeNInt)
		if err != nil {
			logErr(err).Fatalf("Failed to apply %d migrations", relativeNInt)
		}
		logCurrentVersion(mctx, migrate)
		return nil
	},
}

var applyCommand = &cli.Command{
	Name:            "apply",
	Aliases:         []string{"a"},
	Usage:           "Run up migration for specific version",
	ArgsUsage:       "<version>",
	Flags:           MigrateFlags,
	SkipFlagParsing: true,
	Action: func(ctx *cli.Context) error {
		version := ctx.Args().First()
		versionInt, err := strconv.Atoi(version)
		if err != nil {
			logErr(err).Fatal("Unable to parse param <n>")
		}

		log.Infof("Applying version %d", versionInt)

		migrate, mctx, cancel := newMigrateWithCtx(ctx.String("url"), ctx.String("path"))
		defer cancel()
		err = migrate.ApplyVersion(mctx, file.Version(versionInt))
		if err != nil {
			logErr(err).Fatalf("Failed to apply version %d", versionInt)
		}
		logCurrentVersion(mctx, migrate)
		return nil
	},
}

var rollbackCommand = &cli.Command{
	Name:            "rollback",
	Aliases:         []string{"r"},
	Usage:           "Run down migration for specific version",
	ArgsUsage:       "<version>",
	Flags:           MigrateFlags,
	SkipFlagParsing: true,
	Action: func(ctx *cli.Context) error {
		version := ctx.Args().First()
		versionInt, err := strconv.Atoi(version)
		if err != nil {
			logErr(err).Fatal("Unable to parse param <n>")
		}

		log.Infof("Applying version %d", versionInt)

		migrate, mctx, cancel := newMigrateWithCtx(ctx.String("url"), ctx.String("path"))
		defer cancel()
		err = migrate.RollbackVersion(mctx, file.Version(versionInt))
		if err != nil {
			logErr(err).Fatalf("Failed to rollback version %d", versionInt)
		}
		logCurrentVersion(mctx, migrate)
		return nil
	},
}

var gotoCommand = &cli.Command{
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

		migrate, mctx, cancel := newMigrateWithCtx(ctx.String("url"), ctx.String("path"))
		defer cancel()
		currentVersion, err := migrate.Version(mctx)
		if err != nil {
			logErr(err).Fatalf("failed to migrate to version %d", toVersionInt)
		}

		relativeNInt := toVersionInt - int(currentVersion)

		err = migrate.Migrate(mctx, relativeNInt)
		if err != nil {
			logErr(err).Fatalf("Failed to migrate to vefrsion %d", toVersionInt)
		}
		logCurrentVersion(mctx, migrate)
		return nil
	},
}

func newMigrateWithCtx(url, migrationsPath string) (*migrate.Handle, context.Context, func()) {
	done := make(chan struct{})
	m, err := migrate.Open(url, migrationsPath, migrate.WithHooks(
		func(f file.File) error {
			log.Infof("Applying %s migration for version %d (%s)", f.Direction, f.Version, f.Name)
			return nil
		},
		func(f file.File) error {
			done <- struct{}{}
			return nil
		},
	))
	if err != nil {
		log.Fatalf("Initialization failed: %s", err)
	}
	ctx, cancel := newOsInterruptCtx(done)
	return m, ctx, cancel
}

// newOsInterruptCtx returns new context that will be cancelled
// on os.Interrupt signal.
func newOsInterruptCtx(done <-chan struct{}) (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		exit := false
		for loop := true; loop; {
			select {
			case <-done:
				if exit {
					loop = false
				}
			case <-c:
				if exit {
					os.Exit(5)
				}
				cancel()
				exit = true
				log.Info("Aborting after this migration... Hit again to force quit.")
			}
		}
		signal.Stop(c)
	}()
	return ctx, cancel
}

func logErr(err error) *log.Entry {
	return log.WithError(err)
}

func logCurrentVersion(ctx context.Context, migrate *migrate.Handle) {
	version, err := migrate.Version(ctx)
	if err != nil {
		logErr(err).Fatal("Unable to fetch version")
	}
	log.WithField("current-version", version).Info("done")
}
