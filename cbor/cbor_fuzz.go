// +build gofuzz

package cbor

import (
	"bytes"
	"fmt"

	"github.com/polydawn/refmt/shared"
)

func Fuzz(data []byte) int {
	dec := NewDecoder(bytes.NewReader(data))
	buf := &bytes.Buffer{}

	// Run it once to sanitize
	pump := shared.TokenPump{dec, NewEncoder(buf)}
	err := pump.Run()
	if err != nil {
		return 0
	}

	// Run second loop to check stability
	sanitized := buf.Bytes()
	dec = NewDecoder(bytes.NewReader(sanitized))
	buf = &bytes.Buffer{}
	pump = shared.TokenPump{dec, NewEncoder(buf)}
	err = pump.Run()
	if err != nil {
		fmt.Printf("input: %v, sanitized: %v", data, sanitized)
		panic("sanitised failed to prase: " + err.Error())
	}

	out := buf.Bytes()
	if !bytes.Equal(out, sanitized) {
		fmt.Printf("santized: %v, output: %v\n", sanitized, out)
		panic("looping data failed")
	}

	if len(out) == len(data) {
		return 2
	}
	return 1
}
