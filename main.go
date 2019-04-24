package main

import (
	"path"
	"runtime"
	"os"
	"github.com/urfave/cli"
	log "github.com/sirupsen/logrus"	
)
const usage = `zzz`

type ContextHook struct {
}
func (hook ContextHook) Levels() []log.Level {
    return log.AllLevels
}
func (hook ContextHook) Fire(entry *log.Entry) error {
    if pc, file, line, ok := runtime.Caller(8); ok {
        funcName := runtime.FuncForPC(pc).Name()
        entry.Data["file"] = path.Base(file)
        entry.Data["func"] = path.Base(funcName)
        entry.Data["line"] = line
    }
    return nil
}

func main() {
	// http.ListenAndServe("0.0.0.0:8000", nil)

	// log.AddHook(ContextHook{})	
	app := cli.NewApp()
	app.Name = "paddle"
	app.Usage = usage
	

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		stopCommand,
		removeCommand,
		commitCommand,
		listCommand,
		logCommand,
		execCommand,
		networkCommand,
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
