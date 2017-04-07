package journey

import "github.com/urfave/cli"

func Commands() cli.Commands {
	return cli.Commands{
		{
			Name:        "migrate",
			Aliases:     []string{"m"},
			Usage:       "migrate database",
			Subcommands: MigrateCommands(),
		},
		{
			Name:        "schedule",
			Aliases:     []string{"s"},
			Usage:       "Schedule cron jobs",
			Subcommands: CronjobsCommands(),
		},
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "url, u",
			Usage:  "Driver URL",
			Value:  "",
			EnvVar: "DRIVER_URL",
		},
		cli.StringFlag{
			Name:   "path, p",
			Usage:  "Files path",
			Value:  "./files",
			EnvVar: "FILES_PATH",
		},
	}
}
