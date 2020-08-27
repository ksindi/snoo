/*
A CLI to get information from the SNOO Smart Sleeper Bassinet.

		$ snoo sessions --start 2020-08-24 --end 2020-08-24

*/
package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Days prints out daily aggregated session data in CSV format
func Days(c *cli.Context) error {
	username := c.String("username")
	password := c.String("password")

	startTime := c.Timestamp("start")
	endTime := c.Timestamp("end")

	client := NewClient(username, password)
	client.GetHistory(*startTime, *endTime)

	return nil
}

// Sessions prints out session data in CSV format
func Sessions(c *cli.Context) error {
	username := c.String("username")
	password := c.String("password")

	startTime := c.Timestamp("start")
	endTime := c.Timestamp("end")

	client := NewClient(username, password)
	client.GetSessions(*startTime, *endTime)

	return nil
}

// runCli will run the command-line program, see
// https://github.com/urfave/cli
func runCli() error {
	app := cli.NewApp()

	dayFlags := []cli.Flag{
		&cli.TimestampFlag{
			Name:     "start, s",
			Usage:    "start date",
			Layout:   "2006-01-02",
			Required: true,
		},
		&cli.TimestampFlag{
			Name:     "end, e",
			Usage:    "end date",
			Layout:   "2006-01-02",
			Required: true,
		},
	}

	app.Usage = "An API client to the SNOO Smart Sleeper Bassinet"
	app.Commands = []*cli.Command{
		{
			Name:   "sessions",
			Usage:  "print list of sessions for a date range",
			Action: Sessions,
			Flags:  dayFlags,
		},
		{
			Name:   "days",
			Usage:  "print list of daily aggregated sessions for a date range",
			Action: Days,
			Flags:  dayFlags,
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "username, u",
			Usage:    "SNOO login username",
			EnvVars:  []string{"SNOO_USERNAME"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password, p",
			Usage:    "SNOO login password",
			EnvVars:  []string{"SNOO_PASSWORD"},
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable verbose debugging",
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		log.Debugf("Username: '%s'", c.String("username"))

		return nil
	}

	app.EnableBashCompletion = true
	app.Version = Version

	// run the program
	return app.Run(os.Args)
}

func main() {
	// run the cli
	err := runCli()

	// handle the error
	if err != nil {
		log.Fatalln(err)
	}
}
