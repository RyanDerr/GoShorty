package commonflgs

import (
	"github.com/urfave/cli/v2"
)

const (
	ServiceUrlFlag = "service-url"
	serviceUrlEnv  = "GOSHORTY_SERVICE_URL"
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
