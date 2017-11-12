package main

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli"

	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/pretty"
	"github.com/polydawn/refmt/shared"
)

func main() {
	os.Exit(Main(os.Args, os.Stdin, os.Stdout, os.Stderr))
}

func Main(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	app := cli.NewApp()
	app.Name = "refmt"
	app.Authors = []cli.Author{
		cli.Author{Name: "Eric Myhre", Email: "hash@exultant.us"},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Category: "prettyprint",
			Name:     "json=pretty",
			Usage:    "read json, then pretty print it",
			Action: func(c *cli.Context) error {
				return shared.TokenPump{
					json.NewDecoder(stdin),
					pretty.NewEncoder(stdout),
				}.Run()
			},
		},
		cli.Command{
			Category: "prettyprint",
			Name:     "hex=cbor=pretty",
			Usage:    "read cbor in hex, then pretty print it",
			Action: func(c *cli.Context) error {
				return shared.TokenPump{
					cbor.NewDecoder(hexReader(stdin)),
					pretty.NewEncoder(stdout),
				}.Run()
			},
		},
		cli.Command{
			Category: "convert",
			Name:     "json=cbor",
			Usage:    "read json, emit equivalent cbor",
			Action: func(c *cli.Context) error {
				return shared.TokenPump{
					json.NewDecoder(stdin),
					cbor.NewEncoder(stdout),
				}.Run()
			},
		},
		cli.Command{
			Category: "convert",
			Name:     "json=cbor=hex",
			Usage:    "read json, emit equivalent cbor in hex",
			Action: func(c *cli.Context) error {
				return shared.TokenPump{
					json.NewDecoder(stdin),
					cbor.NewEncoder(hexWriter{stdout}),
				}.Run()
			},
		},
		cli.Command{
			Category: "convert",
			Name:     "hex=cbor=json",
			Usage:    "read cbor in hex, emit equivalent json",
			Action: func(c *cli.Context) error {
				return shared.TokenPump{
					cbor.NewDecoder(hexReader(stdin)),
					json.NewEncoder(stdout),
				}.Run()
			},
		},
	}
	app.Writer = stdout
	app.ErrWriter = stderr
	err := app.Run(args)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}
	return 0
}
