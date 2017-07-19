package json

// License note: the string and numeric parsers here borrow
// heavily from the golang stdlib json parser scanner.
// That code is originally Copyright 2010 The Go Authors,
// and is governed by a BSD-style license.

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func (d *Decoder) decodeString() (string, error) {
	// First quote has already been eaten.
	// Start tracking the byte slice; real string starts here.
	d.r.track()
	// Scan until scanner tells us end of string.
	var err error
	for step := strscan_normal; step != nil && err == nil; step, err = step(d.r.readn1()) {
	}
	// Unread one.  The scan loop consumed the trailing quote already,
	// which we don't want to pass onto the parser.
	d.r.unreadn1()
	// Parse!
	s, ok := parseString(d.r.stopTrack())
	if !ok {
		//return string(s), fmt.Errorf("string parse misc fail")
	}
	// Swallow the trailing quote again.
	d.r.readn1()
	return string(s), nil
}

// Scan steps are looped over the stream to find how long the string is.
// A nil step func is returned to indicate the string is done.
// Actually parsing the string is done by 'parseString()'.
type strscanStep func(c byte) (strscanStep, error)

// The default string scanning step state.  Starts here.
func strscan_normal(c byte) (strscanStep, error) {
	if c == '"' { // done!
		return nil, nil
	}
	if c == '\\' {
		return strscan_esc, nil
	}
	if c < 0x20 { // Unprintable bytes are invalid in a json string.
		return nil, fmt.Errorf("invalid unprintable byte in string literal: 0x%x", c)
	}
	return strscan_normal, nil
}

// "esc" is the state after reading `"\` during a quoted string.
func strscan_esc(c byte) (strscanStep, error) {
	switch c {
	case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
		return strscan_normal, nil
	case 'u':
		return strscan_escU, nil
	}
	return nil, fmt.Errorf("invalid byte in string escape sequence: 0x%x", c)
}

// "escU" is the state after reading `"\u` during a quoted string.
func strscan_escU(c byte) (strscanStep, error) {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		return strscan_escU1, nil
	}
	return nil, fmt.Errorf("invalid byte in \\u hexadecimal character escape: 0x%x", c)
}

// "escU1" is the state after reading `"\u1` during a quoted string.
func strscan_escU1(c byte) (strscanStep, error) {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		return strscan_escU12, nil
	}
	return nil, fmt.Errorf("invalid byte in \\u hexadecimal character escape: 0x%x", c)
}

// "escU12" is the state after reading `"\u12` during a quoted string.
func strscan_escU12(c byte) (strscanStep, error) {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		return strscan_escU123, nil
	}
	return nil, fmt.Errorf("invalid byte in \\u hexadecimal character escape: 0x%x", c)
}

// "escU123" is the state after reading `"\u123` during a quoted string.
func strscan_escU123(c byte) (strscanStep, error) {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		return strscan_normal, nil
	}
	return nil, fmt.Errorf("invalid byte in \\u hexadecimal character escape: 0x%x", c)
}

// Convert a json string byte sequence that is a complete string (quotes from
// the outside dropped) bytes ready to be flipped into a go string.
func parseString(s []byte) (t []byte, ok bool) {
	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room?  Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	r, err := strconv.ParseUint(string(s[2:6]), 16, 64)
	if err != nil {
		return -1
	}
	return rune(r)
}

func (d *Decoder) decodeFloat(majorByte byte) (float64, error) {
	// TODO
	return 0, nil
}
