package action

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/CzarSimon/sws/cmd/cli/swsctl/api"
	"github.com/urfave/cli"
)

// Configure gets and stores api configuration information from the user.
func Configure(c *cli.Context) error {
	config := getApiConfig()
	err := config.Save()
	if err != nil {
		fmt.Println("Could not save configuration")
		return err
	}
	return nil
}

// getApiConfig Prompts the user to input api configuration
func getApiConfig() api.Config {
	var config api.Config
	fmt.Println("Enter api configuration")
	config.API.Host = getInput("Host")
	config.API.Port = getInput("Port")
	config.API.Protocol = "http"
	config.Auth.AccessKey = getInput("Access Key")
	return config
}

// getInput Gets user input from stdin
func getInput(key string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(key + ": ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Unable to read the value for '%s'\n", key)
		os.Exit(1)
	}
	return strings.Replace(text, "\n", "", -1)
}
