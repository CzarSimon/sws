package action

import (
	"fmt"
	"os"

	"github.com/CzarSimon/sws/cmd/cli/swsctl/api"
	"github.com/CzarSimon/sws/pkg/service"
	"github.com/urfave/cli"
)

// Apply command for applying a service configuration on a sws node.
func Apply(c *cli.Context) error {
	svcFile := getFirstArg(c, "No service file provided")
	manifest, err := service.ReadService(svcFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	msg, err := api.PostService(manifest.Spec)
	fmt.Print(msg)
	return err
}

// Delete command for deleting a service.
func Delete(c *cli.Context) error {
	var svc service.Service
	svc.Name = getFirstArg(c, "No service name entered")
	msg, err := api.DeleteService(svc)
	fmt.Print(msg)
	return err
}

// List command for listing a users running services.
func List(c *cli.Context) error {
	services, err := api.GetServices()
	if err != nil {
		fmt.Println("Could not get active services")
		return err
	}
	for _, svc := range services {
		fmt.Printf("Name = %s Domain = %s Port = %d Image = %s\n",
			svc.Name, svc.Domain, svc.Port, svc.Image)
	}
	return nil
}

func getFirstArg(c *cli.Context, errMsg string) string {
	val := c.Args().First()
	if val == "" {
		fmt.Println(errMsg)
		os.Exit(1)
	}
	return val
}
