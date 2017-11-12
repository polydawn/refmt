package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"

	"github.com/polydawn/refmt/cbor"
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
			Category: "convert",
			Name:     "hex=cbor=pretty",
			Usage:    "read cbor in hex, then pretty print it",
			Action: func(c *cli.Context) error {
				fmt.Fprintf(stderr, "hello\n")
				in, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				bs := make([]byte, len(in)/2)
				n, err := hex.Decode(bs, in)
				if err != nil {
					return err
				}
				if n != len(in)/2 {
					return fmt.Errorf("hex len mismatch: %d chars became %d bytes", len(in), n)
				}

				pump := shared.TokenPump{
					cbor.NewDecoder(bytes.NewBuffer(bs)),
					pretty.NewEncoder(os.Stdout),
				}
				return pump.Run()
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
