package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/pretty"
	"github.com/polydawn/refmt/shared"
)

func main() {
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	bs := make([]byte, len(in)/2)
	n, err := hex.Decode(bs, in)
	if err != nil {
		panic(err)
	}
	if n != len(in)/2 {
		panic(fmt.Errorf("hex len mismatch: %d chars became %d bytes", len(in), n))
	}

	pump := shared.TokenPump{
		cbor.NewDecoder(bytes.NewBuffer(bs)),
		pretty.NewEncoder(os.Stdout),
	}
	if err := pump.Run(); err != nil {
		panic(err)
	}
}
