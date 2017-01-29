package json

import (
	"fmt"
	"io"
	"strconv"

	. "github.com/polydawn/go-xlate/tok"
)

func NewSerializer(wr io.Writer) *Serializer {
	return &Serializer{wr: wr}
}

/*
	A json.Serializer is a TokenSink implementation that emits json bytes.
*/
type Serializer struct {
	wr io.Writer

	// Stack, tracking how many array and map opens are outstanding.
	// (Values are only 'phase_mapExpectKeyOrEnd' and 'phase_arrExpectValueOrEnd'.)
	stack   []phase
	current phase // shortcut to value at end of stack
	some    bool  // set to true after first value in any context; use to append commas.

	// Spare memory, for use in operations on leaf nodes (e.g. temp space for an int serialization).
	scratch [64]byte
}

type phase int

const (
	phase_anyExpectValue phase = iota
	phase_mapExpectKeyOrEnd
	phase_mapExpectValue
	phase_arrExpectValueOrEnd
)

func (d *Serializer) Step(tok *Token) (done bool, err error) {
	switch d.current {
	case phase_anyExpectValue:
		switch *tok {
		case Token_MapOpen:
			d.pushPhase(phase_mapExpectKeyOrEnd)
			d.wr.Write(wordMapOpen)
			return false, nil
		case Token_ArrOpen:
			d.pushPhase(phase_arrExpectValueOrEnd)
			d.wr.Write(wordArrOpen)
			return false, nil
		case Token_MapClose:
			return true, fmt.Errorf("unexpected mapClose; expected start of value")
		case Token_ArrClose:
			return true, fmt.Errorf("unexpected arrClose; expected start of value")
		default:
			// It's a value; handle it.
			d.flushValue(tok)
			return true, nil
		}
	case phase_mapExpectKeyOrEnd:
		switch *tok {
		case Token_MapOpen:
			return true, fmt.Errorf("unexpected mapOpen; expected start of key or end of map")
		case Token_ArrOpen:
			return true, fmt.Errorf("unexpected arrOpen; expected start of key or end of map")
		case Token_MapClose:
			d.wr.Write(wordMapClose)
			return d.popPhase()
		case Token_ArrClose:
			return true, fmt.Errorf("unexpected arrClose; expected start of key or end of map")
		default:
			// It's a key.  It'd better be a string.
			switch v2 := (*tok).(type) {
			case *string:
				d.entrySep()
				d.wr.Write([]byte(fmt.Sprintf("%q", *v2)))
				d.wr.Write(wordColon)
				d.current = phase_mapExpectValue
				return false, nil
			default:
				return true, fmt.Errorf("unexpected token of type %T; expected map key", *tok)
			}
		}
	case phase_mapExpectValue:
		switch *tok {
		case Token_MapOpen:
			d.pushPhase(phase_mapExpectKeyOrEnd)
			d.wr.Write(wordMapOpen)
			return false, nil
		case Token_ArrOpen:
			d.pushPhase(phase_arrExpectValueOrEnd)
			d.wr.Write(wordArrOpen)
			return false, nil
		case Token_MapClose:
			return true, fmt.Errorf("unexpected mapClose; expected start of value")
		case Token_ArrClose:
			return true, fmt.Errorf("unexpected arrClose; expected start of value")
		default:
			// It's a value; handle it.
			d.flushValue(tok)
			d.current = phase_mapExpectKeyOrEnd
			return false, nil
		}
	case phase_arrExpectValueOrEnd:
		switch *tok {
		case Token_MapOpen:
			d.entrySep()
			d.pushPhase(phase_mapExpectKeyOrEnd)
			d.wr.Write(wordMapOpen)
			return false, nil
		case Token_ArrOpen:
			d.entrySep()
			d.pushPhase(phase_arrExpectValueOrEnd)
			d.wr.Write(wordArrOpen)
			return false, nil
		case Token_MapClose:
			return true, fmt.Errorf("unexpected mapClose; expected start of value or end of array")
		case Token_ArrClose:
			d.wr.Write(wordArrClose)
			return d.popPhase()
		default:
			// It's a value; handle it.
			d.entrySep()
			d.flushValue(tok)
			return false, nil
		}
	default:
		panic("Unreachable")
	}
}

func (d *Serializer) pushPhase(p phase) {
	d.current = p
	d.stack = append(d.stack, d.current)
	d.some = false
}

// Pop a phase from the stack; return 'true' if stack now empty.
func (d *Serializer) popPhase() (bool, error) {
	n := len(d.stack) - 1
	if n == 0 {
		return true, nil
	}
	if n < 0 { // the state machines are supposed to have already errored better
		panic("jsonSerializer stack overpopped")
	}
	d.current = d.stack[n-1]
	d.stack = d.stack[0:n]
	d.some = true
	return false, nil
}

// The most heavily used words, cached as byte slices.
var (
	wordTrue     = []byte("true")
	wordFalse    = []byte("false")
	wordNull     = []byte("null")
	wordArrOpen  = []byte("[")
	wordArrClose = []byte("]")
	wordMapOpen  = []byte("{")
	wordMapClose = []byte("}")
	wordColon    = []byte(":")
	wordComma    = []byte(",")
)

// Emit an entry separater (comma), unless we're at the start of an object.
// Mark that we *do* have some content, regardless, so next time will need a sep.
func (d *Serializer) entrySep() {
	if d.some {
		d.wr.Write(wordComma)
	}
	d.some = true
}

func (d *Serializer) flushValue(tokSlot *Token) {
	switch valp := (*tokSlot).(type) {
	case *string:
		d.wr.Write([]byte(fmt.Sprintf("%q", *valp)))
	case *int:
		b := strconv.AppendInt(d.scratch[:0], int64(*valp), 10)
		d.wr.Write(b)
	case nil:
		d.wr.Write(wordNull)
	default:
		panic(fmt.Errorf("TODO finish more jsonSerializer primitives support: type %T", *tokSlot))
	}
}
