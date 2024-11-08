package shortencmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RyanDerr/GoShorty/api/modules/shorten"
	"github.com/RyanDerr/GoShorty/internal/cmd/flags/commonflgs"
	"github.com/RyanDerr/GoShorty/internal/cmd/flags/shortenflgs"
	"github.com/urfave/cli/v2"
)

const (
	ShortenCommandName = "shorten"
)

var Command = &cli.Command{
	Name:  ShortenCommandName,
	Usage: "Shorten a URL",
	Flags: shortenflgs.GetShortenFlags(),
	Action: func(ctx *cli.Context) error {
		serviceUrl, err := commonflgs.GetServiceURL(ctx)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		ShortenCommandParams := &shorten.ShortenRequest{
			URL:         ctx.String(shortenflgs.URLFlag),
			CustomShort: ctx.String(shortenflgs.ShortFlag),
			Expiration:  ctx.String(shortenflgs.TtlFlag),
		}

		// Call domain service to handle the URL shortening
		resp, err := shortenURL(ShortenCommandParams, serviceUrl)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		fmt.Printf("%v\n", resp.String())
		return nil
	},
}

func shortenURL(data *shorten.ShortenRequest, serviceUrl string) (*shorten.ShortenResponse, error) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %v", err)
	}

	resp, err := http.Post(serviceUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned non-201 status: %v", resp.Status)
	}

	var response shorten.ShortenResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}
