package main

import (
	"os"
	"github.com/urfave/cli"
	log "github.com/sirupsen/logrus"	
)
const usage = `zzz`

func main() {
	app := cli.NewApp()
	app.Name = "paddle"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})

		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}	
}
