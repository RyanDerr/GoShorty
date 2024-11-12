package shortencmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RyanDerr/GoShorty/internal/cmd/flags/shortenflgs"
	"github.com/RyanDerr/GoShorty/pkg/request"
	"github.com/RyanDerr/GoShorty/pkg/response"
	"github.com/urfave/cli/v2"
)

const (
	ShortenCommandName = "shorten"
)

var apiResponse struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

var Command = &cli.Command{
	Name:  ShortenCommandName,
	Usage: "Shorten a URL",
	Flags: shortenflgs.GetShortenFlags(),
	Action: func(ctx *cli.Context) error {
		serviceUrl, err := shortenflgs.GetShortenServiceUrl(ctx)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		ShortenCommandParams := &request.ShortenUrlRequest{
			Url:         ctx.String(shortenflgs.URLFlag),
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

func shortenURL(data *request.ShortenUrlRequest, serviceUrl string) (*response.ShortenUrlResponse, error) {

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
		return nil, fmt.Errorf("API returned status: %v", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var response response.ShortenUrlResponse
	if err := json.Unmarshal(apiResponse.Data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return &response, nil
}
