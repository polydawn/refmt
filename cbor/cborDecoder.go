package cbor

import (
	"fmt"
	"io"

	. "github.com/polydawn/go-xlate/tok"
)

type Decoder struct {
	r quickReader

	stack []decoderStep // When empty, and step returns done, all done.
	step  decoderStep   // Shortcut to end of stack.
	left  []int         // Statekeeping space for definite-len map and array.

	spareBytes []byte
}

func NewDecoder(r io.Reader) (d *Decoder) {
	d = &Decoder{
		r:     &quickReaderStream{br: &readerByteScanner{r: r}},
		stack: make([]decoderStep, 0, 10),
		left:  make([]int, 0, 10),
	}
	d.step = d.step_acceptValue
	return
}

func (d *Decoder) Reset() {
	d.stack = d.stack[0:0]
	d.step = d.step_acceptValue
	d.left = d.left[0:0]
}

type decoderStep func(tokenSlot *Token) (done bool, err error)

func (d *Decoder) Step(tokenSlot *Token) (done bool, err error) {
	done, err = d.step(tokenSlot)
	// If the step errored: out, entirely.
	if err != nil {
		return true, err
	}
	// If the step wasn't done, return same status.
	if !done {
		return false, nil
	}
	// If it WAS done, pop next, or if stack empty, we're entirely done.
	nSteps := len(d.stack) - 1
	if nSteps <= 0 {
		return true, nil // that's all folks
	}
	d.step = d.stack[nSteps]
	d.stack = d.stack[0:nSteps]
	return false, nil
}

func (d *Decoder) pushPhase(newPhase decoderStep) {
	d.stack = append(d.stack, d.step)
	d.step = newPhase
}

// The original step, where any value is accepted, and no terminators for composites are valid.
// ONLY used in the original step; all other steps handle leaf nodes internally.
func (d *Decoder) step_acceptValue(tokenSlot *Token) (done bool, err error) {
	majorByte := d.r.readn1()
	return d.stepHelper_acceptValue(majorByte, tokenSlot)
}

// Step in midst of decoding an indefinite-length array.
func (d *Decoder) step_acceptArrValueOrBreak(tokenSlot *Token) (done bool, err error) {
	majorByte := d.r.readn1()
	switch majorByte {
	case cborSigilBreak:
		*tokenSlot = Token_ArrClose
		return true, nil
	default:
		_, err := d.stepHelper_acceptValue(majorByte, tokenSlot)
		return false, err
	}
}

// Step in midst of decoding an indefinite-length map, key expected up next, or end.
func (d *Decoder) step_acceptMapIndefKey(tokenSlot *Token) (done bool, err error) {
	majorByte := d.r.readn1()
	switch majorByte {
	case cborSigilBreak:
		*tokenSlot = Token_MapClose
		return true, nil
	default:
		d.step = d.step_acceptMapIndefValueOrBreak
		_, err := d.stepHelper_acceptValue(majorByte, tokenSlot) // FIXME surely not *any* value?  not composites, at least?
		return false, err
	}
}

// Step in midst of decoding an indefinite-length map, value expected up next.
func (d *Decoder) step_acceptMapIndefValueOrBreak(tokenSlot *Token) (done bool, err error) {
	majorByte := d.r.readn1()
	switch majorByte {
	case cborSigilBreak:
		return true, fmt.Errorf("unexpected break; expected value in indefinite-length map")
	default:
		d.step = d.step_acceptMapIndefKey
		_, err = d.stepHelper_acceptValue(majorByte, tokenSlot)
		return false, err
	}
}

// Step in midst of decoding a definite-length array.
func (d *Decoder) step_acceptArrValue(tokenSlot *Token) (done bool, err error) {
	// Yield close token and return done flag if expecting no more entries.
	ll := len(d.left) - 1
	if d.left[ll] == 0 {
		d.left = d.left[0:ll]
		*tokenSlot = Token_ArrClose
		return true, nil
	}
	d.left[ll]--
	// Read next value.
	majorByte := d.r.readn1()
	_, err = d.stepHelper_acceptValue(majorByte, tokenSlot)
	return false, err
}

// Step in midst of decoding an definite-length map, key expected up next.
func (d *Decoder) step_acceptMapKey(tokenSlot *Token) (done bool, err error) {
	// Read next key.
	majorByte := d.r.readn1()
	d.step = d.step_acceptMapValue
	_, err = d.stepHelper_acceptValue(majorByte, tokenSlot) // FIXME surely not *any* value?  not composites, at least?
	return false, err
}

// Step in midst of decoding an definite-length map, value expected up next.
func (d *Decoder) step_acceptMapValue(tokenSlot *Token) (done bool, err error) {
	// Read next value.
	majorByte := d.r.readn1()
	_, err = d.stepHelper_acceptValue(majorByte, tokenSlot)
	// If expecting no more entries, pop state
	// and set next step to endMap instead of acceptKey.
	ll := len(d.left) - 1
	if d.left[ll] <= 1 {
		d.left = d.left[0:ll]
		d.step = d.step_endMap
		return false, err
	}
	d.left[ll]--
	d.step = d.step_acceptMapKey
	return false, err
}

// Step when reached the expected end of a definite-length map.
func (d *Decoder) step_endMap(tokenSlot *Token) (done bool, err error) {
	*tokenSlot = Token_MapClose
	return true, nil
}

func (d *Decoder) stepHelper_acceptValue(majorByte byte, tokenSlot *Token) (done bool, err error) {
	switch majorByte {
	case cborSigilNil:
		*tokenSlot = nil
		return true, nil
	case cborSigilFalse:
		*tokenSlot = false
		return true, nil
	case cborSigilTrue:
		*tokenSlot = true
		return true, nil
	case cborSigilFloat16, cborSigilFloat32, cborSigilFloat64:
		*tokenSlot = d.decodeFloat(majorByte)
		return true, nil
	case cborSigilIndefiniteBytes:
		var x []byte
		x, err = d.decodeBytesIndefinite(nil)
		*tokenSlot = &x
		return true, err
	case cborSigilIndefiniteString:
		var x string
		x, err = d.decodeStringIndefinite()
		*tokenSlot = &x
		return true, err
	case cborSigilIndefiniteArray:
		*tokenSlot = Token_ArrOpen
		d.pushPhase(d.step_acceptArrValueOrBreak)
		return false, nil
	case cborSigilIndefiniteMap:
		*tokenSlot = Token_MapOpen
		d.pushPhase(d.step_acceptMapIndefKey)
		return false, nil
	default:
		switch {
		case majorByte >= cborMajorUint && majorByte < cborMajorNegInt:
			var x uint64
			x, err = d.decodeUint(majorByte)
			*tokenSlot = &x
			return true, err
		case majorByte >= cborMajorNegInt && majorByte < cborMajorBytes:
			var x int64
			x, err = d.decodeNegInt(majorByte)
			*tokenSlot = &x
			return true, err
		case majorByte >= cborMajorBytes && majorByte < cborMajorString:
			var x []byte
			x, err = d.decodeBytes(majorByte)
			*tokenSlot = &x
			return true, err
		case majorByte >= cborMajorString && majorByte < cborMajorArray:
			var x string
			x, err = d.decodeString(majorByte)
			*tokenSlot = &x
			return true, err
		case majorByte >= cborMajorArray && majorByte < cborMajorMap:
			*tokenSlot = Token_ArrOpen
			d.pushPhase(d.step_acceptArrValue)
			var n int
			n, err = d.decodeLen(majorByte)
			d.left = append(d.left, n)
			return false, nil
		case majorByte >= cborMajorMap && majorByte < cborMajorTag:
			*tokenSlot = Token_MapOpen
			var n int
			n, err = d.decodeLen(majorByte)
			if err != nil {
				return true, err
			}
			if n == 0 {
				d.pushPhase(d.step_endMap)
				return false, nil
			}
			d.left = append(d.left, n)
			d.pushPhase(d.step_acceptMapKey)
			return false, nil
		case majorByte >= cborMajorTag && majorByte < cborMajorSimple:
			return true, fmt.Errorf("cbor tags not supported")
			// but when we do, it'll be by saving it as another field of the Token, and recursing.
		default:
			return true, fmt.Errorf("Invalid majorByte: 0x%x", majorByte)
		}
	}
}
