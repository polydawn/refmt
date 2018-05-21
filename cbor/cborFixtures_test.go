package cbor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"

	. "github.com/warpfork/go-wish"

	. "github.com/polydawn/refmt/tok"
	"github.com/polydawn/refmt/tok/fixtures"
)

func Test(t *testing.T) {
	testBool(t)
	testString(t)
	testMap(t)
	testArray(t)
	testComposite(t)
}

func checkEncoding(t *testing.T, sequence fixtures.Sequence, expectSerial []byte, expectErr error) {
	t.Helper()
	outputBuf := &bytes.Buffer{}
	tokenSink := NewEncoder(outputBuf)

	// Run steps, advancing through the token sequence.
	//  If it stops early, just report how many steps in; we Wish on that value.
	//  If it doesn't stop in time, just report that bool; we Wish on that value.
	var nStep int
	var done bool
	var err error
	for _, tok := range sequence.Tokens {
		nStep++
		done, err = tokenSink.Step(&tok)
		if done || err != nil {
			break
		}
	}

	// Assert final result.
	Wish(t, done, ShouldEqual, true)
	Wish(t, nStep, ShouldEqual, len(sequence.Tokens))
	Wish(t, err, ShouldEqual, expectErr)
	Wish(t, outputBuf.Bytes(), ShouldEqual, expectSerial)
}

func checkDecoding(t *testing.T, expectSequence fixtures.Sequence, serial []byte, expectErr error) {
	t.Helper()
	inputBuf := bytes.NewBuffer(serial)
	tokenSrc := NewDecoder(inputBuf)

	// Run steps, advancing until the decoder reports it's done.
	//  If the decoder keeps yielding more tokens than we expect, that's fine...
	//  we just keep recording them, and we'll diff later.
	//  There's a cutoff when it overshoots by 10 tokens because generally
	//  that indicates we've found some sort of loop bug and 10 extra token
	//  yields is typically enough info to diagnose with.
	var nStep int
	var done bool
	var yield = make([]Token, len(expectSequence.Tokens)+10)
	var err error
	for ; nStep <= len(expectSequence.Tokens)+10; nStep++ {
		done, err = tokenSrc.Step(&yield[nStep])
		if done || err != nil {
			break
		}
	}
	nStep++
	yield = yield[:nStep]

	// Assert final result.
	Wish(t, done, ShouldEqual, true)
	Wish(t, nStep, ShouldEqual, len(expectSequence.Tokens))
	Wish(t, yield, ShouldEqual, expectSequence.Tokens)
	Wish(t, err, ShouldEqual, expectErr)
}

func bcat(bss ...[]byte) []byte {
	l := 0
	for _, bs := range bss {
		l += len(bs)
	}
	rbs := make([]byte, 0, l)
	for _, bs := range bss {
		rbs = append(rbs, bs...)
	}
	return rbs
}

func b(b byte) []byte { return []byte{b} }

func deB64(s string) []byte {
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return bs
}

var inapplicable = fmt.Errorf("skipme: inapplicable")

var cborFixtures = []struct {
	title        string
	sequence     fixtures.Sequence
	serial       []byte
	encodeResult error
	decodeResult error
}{
	// Numbers.
	{"",
		fixtures.Sequence{"integer zero", []Token{{Type: TInt, Int: 0}}},
		deB64("AA=="),
		nil,
		inapplicable, // Impossible to decode to this token because cbor doens't disambiguate positive vs signed ints.
	},
	{"",
		fixtures.Sequence{"integer zero unsigned", []Token{{Type: TUint, Uint: 0}}},
		deB64("AA=="),
		nil,
		nil,
	},
	{"",
		fixtures.Sequence{"integer one", []Token{{Type: TInt, Int: 1}}},
		deB64("AQ=="),
		nil,
		inapplicable, // Impossible to decode to this token because cbor doens't disambiguate positive vs signed ints.
	},
	{"",
		fixtures.Sequence{"integer one unsigned", []Token{{Type: TUint, Uint: 1}}},
		deB64("AQ=="),
		nil,
		nil,
	},
	{"",
		fixtures.Sequence{"integer neg 1", []Token{{Type: TInt, Int: -1}}},
		deB64("IA=="),
		nil,
		nil,
	},
	{"",
		fixtures.Sequence{"integer neg 100", []Token{{Type: TInt, Int: -100}}},
		deB64("OGM="),
		nil,
		nil,
	},
	{"",
		fixtures.Sequence{"integer 1000000", []Token{{Type: TInt, Int: 1000000}}},
		deB64("GgAPQkA="),
		nil,
		inapplicable, // Impossible to decode to this token because cbor doens't disambiguate positive vs signed ints.
	},
	{"",
		fixtures.Sequence{"integer 1000000 unsigned", []Token{{Type: TUint, Uint: 1000000}}},
		deB64("GgAPQkA="),
		nil,
		nil,
	},
	//	{"",  // This fixture expects the float32 encoding, and we currently lack support for detecting when things can be safely packed thusly.
	//		fixtures.Sequence{"float decimal e+38", []Token{{Type: TFloat64, Float64: 3.4028234663852886e+38}}},
	//		deB64("+n9///8="),
	//		nil,nil,
	//	},
	{"",
		fixtures.Sequence{"float 1 e+100", []Token{{Type: TFloat64, Float64: 1.0e+300}}},
		deB64("+3435DyIAHWc"),
		nil,
		nil,
	},

	// Byte strings.
	{"",
		fixtures.SequenceMap["short byte array"],
		bcat(b(0x40+5), []byte(`value`)),
		nil,
		nil,
	},
	{"indefinite length bytes (single actual hunk)",
		fixtures.SequenceMap["short byte array"],
		bcat(b(0x5f), b(0x40+5), []byte(`value`), b(0xff)),
		inapplicable,
		nil,
	},
	{"indefinite length bytes (multiple hunks)",
		fixtures.SequenceMap["short byte array"],
		bcat(b(0x5f), b(0x40+2), []byte(`va`), b(0x40+3), []byte(`lue`), b(0xff)),
		inapplicable,
		nil,
	},
	{"",
		fixtures.SequenceMap["long zero byte array"],
		bcat(b(0x40+0x19), []byte{0x1, 0x90}, bytes.Repeat(b(0x0), 400)),
		nil,
		nil,
	},

	// Tags.
	{"",
		fixtures.SequenceMap["tagged object"],
		bcat(b(0xc0+(0x20-8)), b(50), b(0xa0+1), b(0x60+1), []byte(`k`), b(0x60+1), []byte(`v`)),
		nil,
		nil,
	},
	{"",
		fixtures.SequenceMap["tagged string"],
		bcat(b(0xc0+(0x20-8)), b(50), b(0x60+5), []byte(`wahoo`)),
		nil,
		nil,
	},
	{"",
		fixtures.SequenceMap["array with mixed tagged values"],
		bcat(b(0x80+2),
			b(0xc0+(0x20-8)), b(40), b(0x00+(0x19)), []byte{0x1, 0x90},
			b(0xc0+(0x20-8)), b(50), b(0x60+3), []byte(`500`)),
		nil,
		nil,
	},
	{"",
		fixtures.SequenceMap["object with deeper tagged values"],
		bcat(b(0xa0+5),
			b(0x60+2), []byte(`k1`), b(0xc0+(0x20-8)), b(50), b(0x60+3), []byte(`500`),
			b(0x60+2), []byte(`k2`), b(0x60+8), []byte(`untagged`),
			b(0x60+2), []byte(`k3`), b(0xc0+(0x20-8)), b(60), b(0x60+3), []byte(`600`),
			b(0x60+2), []byte(`k4`), b(0x80+2),
			/**/ b(0xc0+(0x20-8)), b(50), b(0x60+4), []byte(`asdf`),
			/**/ b(0xc0+(0x20-8)), b(50), b(0x60+4), []byte(`qwer`),
			b(0x60+2), []byte(`k5`), b(0xc0+(0x20-8)), b(50), b(0x60+3), []byte(`505`),
		),
		nil,
		nil,
	},
}
