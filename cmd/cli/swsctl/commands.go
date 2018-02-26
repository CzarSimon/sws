package main

import (
	"fmt"

	"github.com/CzarSimon/sws/cmd/cli/swsctl/action"
	"github.com/urfave/cli"
)

// ConfigureCommand configures swsctl for use.
func ConfigureCommand() cli.Command {
	return cli.Command{
		Name:   "configure",
		Usage:  fmt.Sprintf("Configures %s for use", APP_NAME),
		Action: action.Configure,
	}
}

// ApplyCommand applies a service specification on the remote sws node.
func ApplyCommand() cli.Command {
	return cli.Command{
		Name:   "apply",
		Usage:  "Applies a service specification on the remote sws node",
		Action: action.Apply,
	}
}

// DeleteCommand deletes a service from the remote sws node.
func DeleteCommand() cli.Command {
	return cli.Command{
		Name:   "delete",
		Usage:  "Deletes a service from the remote sws node",
		Action: action.Delete,
	}
}

// ListCommand lists a users running services.
func ListCommand() cli.Command {
	return cli.Command{
		Name:   "ls",
		Usage:  "Lists your current running services",
		Action: action.List,
	}
}
