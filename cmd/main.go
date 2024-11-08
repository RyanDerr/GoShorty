package main

import (
	"log"
	"os"

	"github.com/RyanDerr/GoShorty/internal/cmd/commands/shortencmd"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "goshorty",
		Usage: "CLI for GoShorty URL shortener",
		Commands: []*cli.Command{
			shortencmd.Command,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
