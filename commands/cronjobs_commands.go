package journey

import (
	"github.com/Sirupsen/logrus"
	"github.com/db-journey/cronjobs"
	"github.com/db-journey/migrate/driver"
	"github.com/urfave/cli"
)

var CronjobsFlags = []cli.Flag{}

func CronjobsCommands() cli.Commands {
	return cli.Commands{
		startCommand,
	}
}

var startCommand = cli.Command{
	Name:    "start",
	Aliases: []string{"s"},
	Usage:   "Start scheduler",
	Flags:   CronjobsFlags,
	Action: func(ctx *cli.Context) error {

		driver, err := driver.New(ctx.GlobalString("url"))
		if err != nil {
			logrus.WithError(err).Fatal("Can't initiate driver")
		}

		scheduler := cronjobs.New(driver)
		logrus.Info("Loading cron files from ", ctx.GlobalString("path"))
		err = scheduler.ReadFiles(ctx.GlobalString("path"))
		if err != nil {
			logrus.WithError(err).Fatal("Can't load files")
		}
		scheduler.Logger = func(runs chan *cronjobs.Run) {
			for run := range runs {
				logger := logrus.WithField("name", run.Name)
				if run.Error != nil {
					logger.WithError(run.Error).Error("Failed to run job")
					continue
				}
				logger.Info("Running")
			}
		}
		numJobs := len(scheduler.Entries())
		if numJobs == 0 {
			logrus.Fatal("No cron job found in ", ctx.GlobalString("path"))
		}
		logrus.Infof("%d job(s) loaded", numJobs)
		logrus.Info("Starting Scheduler")
		scheduler.Start()
		select {}
		return nil
	},
}
