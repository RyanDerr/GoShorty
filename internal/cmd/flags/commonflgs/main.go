package commonflgs

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const (
	ServiceUrlFlag = "service-url"
	serviceUrlEnv  = "GOSHORTY_SERVICE_URL"
	apiEndpoint    = "/api/v1"
)

var commonFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    ServiceUrlFlag,
		Usage:   "URL of the GoShorty service",
		EnvVars: []string{serviceUrlEnv},
	},
}

func GetCommonFlags() []cli.Flag {
	return commonFlags
}

func GetServiceURL(c *cli.Context) (string, error) {
	serviceURL := c.String(ServiceUrlFlag)
	if serviceURL == "" {
		return "", fmt.Errorf("service URL not provided for the GoShorty service")
	}
	return fmt.Sprintf("%s%s", serviceURL, apiEndpoint), nil
}
