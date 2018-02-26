package main

import (
	"os"

	"github.com/urfave/cli"
)

const (
	APP_NAME    = "swsctl"
	APP_USAGE   = "cli tool for managing services on sws"
	APP_VERSION = "0.4"
)

// getApp Sets up a cli app and returns it
func getApp() *cli.App {
	app := cli.NewApp()
	app.Name = APP_NAME
	app.Usage = APP_USAGE
	app.Version = APP_VERSION
	app.Commands = getAppCommands()
	return app
}

// getAppCommands Returns a list of cli subcommands
func getAppCommands() []cli.Command {
	return []cli.Command{
		ConfigureCommand(),
		ApplyCommand(),
		DeleteCommand(),
		ListCommand(),
	}
}

func main() {
	app := getApp()
	app.Run(os.Args)
}
