package shortenflgs

import (
	"fmt"

	api "github.com/RyanDerr/GoShorty/api/routes"
	"github.com/RyanDerr/GoShorty/internal/cmd/flags/commonflgs"
	"github.com/urfave/cli/v2"
)

const (
	URLFlag   = "url"
	ShortFlag = "short"
	TtlFlag   = "ttl"
)

var shortenFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     URLFlag,
		Aliases:  []string{"u"},
		Usage:    "URL to shorten",
		Required: true,
	},
	&cli.StringFlag{
		Name:     ShortFlag,
		Aliases:  []string{"s"},
		Usage:    "Custom short URL",
		Required: false,
	},
	&cli.StringFlag{
		Name:     TtlFlag,
		Aliases:  []string{"t"},
		Usage:    "Time to live for the shortened URL",
		Required: false,
	},
}

func GetShortenFlags() []cli.Flag {
	shortenFlags = append(shortenFlags, commonflgs.GetCommonFlags()...)
	return shortenFlags
}

func GetShortenServiceUrl(c *cli.Context) (string, error) {
	serviceURL := c.String(commonflgs.ServiceUrlFlag)
	if serviceURL == "" {
		return "", fmt.Errorf("service URL not provided for the GoShorty service")
	}
	return fmt.Sprintf("%s%s", serviceURL, api.GetShortenRoute()), nil
}
