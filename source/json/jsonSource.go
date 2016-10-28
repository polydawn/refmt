package json

import (
	"bytes"
	"io"

	"github.com/polydawn/go-xlate"
)

type reader struct {
	r   io.Reader
	one []byte
	err error
}

func newReader(r io.Reader) *reader {
	return &reader{
		r:   r,
		one: make([]byte, 1),
	}
}

// Read a byte.  Stores error in reader struct.
// Continuing to call ReadByte without checking for errors may yield the same byte in subsequent calls.
func (r *reader) ReadByte() byte {
	_, r.err = r.r.Read(r.one)
	return r.one[0]
}

// decodeState represents the state while decoding a JSON value.
type decodeState struct {
	// This is *substantially* different than the stdlib impl.
	// The stdlib json impl runs the scanner to completion just to find any value end-to-end,
	// and buffers all those bytes entirely (in a `data []byte` field in the comparable struct).
	// We don't really want to run the scanner twice or think that's a reasonable
	// way to buffer, given that our purpose is to emit another token *stream*;
	// so, we keep a streaming byte reader through this far.

	bytes     *reader
	tokenSink xlate.Destination

	scan scanner
}

// scanWhile processes bytes from the reader until it
// receives a scan code not equal to op, then returns the new scan code.
func (d *decodeState) scanWhile(op scanTransition) scanTransition {
	var newOp scanTransition
	for {
		newOp = d.scan.step(&d.scan, d.bytes.ReadByte())
		if newOp != op {
			break
		}
	}
	return newOp
}

func (d *decodeState) emitValue() {
	switch op := d.scanWhile(scanTransition_SkipSpace); op {
	case scanTransition_BeginArray:
		d.emitArray()
	case scanTransition_BeginMap:
		d.emitMap()
	case scanTransition_BeginLiteral:
		d.emitLiteral()
	case scanTransition_Error:
		panic(d.scan.err)
	default:
		panic("invalid state")
	}
}

func (d *decodeState) emitArray() {
	d.tokenSink.OpenArray()
	for {
		op := d.scanWhile(scanTransition_SkipSpace)
		switch op {
		case scanTransition_EndArray:
			d.tokenSink.CloseArray()
			return
		case scanTransition_BeginLiteral:
			d.emitValue()
		default:
			panic("invalid state")
		}
	}
}

func (d *decodeState) emitMap() {
	d.tokenSink.OpenMap()
	for {
		// either the map ends, or we should consume a new literal for a key.
		op := d.scanWhile(scanTransition_SkipSpace)
		switch op {
		case scanTransition_EndMap:
			d.tokenSink.CloseMap()
			return
		case scanTransition_BeginLiteral:
			// TODO consume the key
			d.tokenSink.WriteMapKey("todo")
		default:
			panic("invalid state")
		}
		// a value must follow.
		d.emitValue()
	}
}

func (d *decodeState) emitLiteral() {
	// All bytes inside literal return scanContinue op code.
	// TODO check: there is a terrifying possibility we need to conjure the *prev* byte here
	//  Yep looks like it: the leading quote on a string is gone, and that's fine, but not so much for numbers.
	var newOp scanTransition
	var byt byte
	buf := &bytes.Buffer{} // TODO this very often does *not* need a heap alloc...
	buf.Write(d.bytes.one) // take the last byte read; scanner is always overstepping by one.
	for {
		byt = d.bytes.ReadByte()
		newOp = d.scan.step(&d.scan, byt)
		if newOp != scanTransition_Continue {
			break
		}
		buf.WriteByte(byt)
	}
	s := buf.String()
	switch s {
	//case isNumber: // FIXME YEAH PRETTY SURE SCANNER SHOULD SAY SO
	case "null":
		d.tokenSink.WriteNull()
	case "false":
		panic("todo")
	case "true":
		panic("todo")
	default:
		d.tokenSink.WriteString(s)
	}
}

func (d *jsonDecDriver) appendStringAsBytes() {
	if d.tok == 0 {
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() {
		}
		d.tok = b
	}
	if d.tok != '"' {
		d.d.errorf("json: expect char '%c' but got char '%c'", '"', d.tok)
	}
	d.tok = 0

	v := d.bs[:0]
	var c uint8
	r := d.r
	for {
		c = r.readn1()
		if c == '"' {
			break
		} else if c == '\\' {
			c = r.readn1()
			switch c {
			case '"', '\\', '/', '\'':
				v = append(v, c)
			case 'b':
				v = append(v, '\b')
			case 'f':
				v = append(v, '\f')
			case 'n':
				v = append(v, '\n')
			case 'r':
				v = append(v, '\r')
			case 't':
				v = append(v, '\t')
			case 'u':
				rr := d.jsonU4(false)
				if utf16.IsSurrogate(rr) {
					rr = utf16.DecodeRune(rr, d.jsonU4(true))
				}
				w2 := utf8.EncodeRune(d.bstr[:], rr)
				v = append(v, d.bstr[:w2]...)
			default:
				d.d.errorf("json: unsupported escaped value: %c", c)
			}
		} else {
			v = append(v, c)
		}
	}
	d.bs = v
}
