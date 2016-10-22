package json

import (
	"fmt"
	"io"

	"github.com/polydawn/go-xlate"
)

/*
	Returns a new `Destination` implementation that will write the value(s)
	it receives into a the given `interface{}` reference.
*/
func New(w io.Writer) xlate.Destination {
	return &dest{w: w}
}

type dest struct {
	w io.Writer

	// stack o opens
	// trailing bit for if not first elem
	// recycles some magic words because i felt like it
	ctxStack [][]byte
}

func (d *dest) pushCtx(word []byte) {
	d.ctxStack = append(d.ctxStack, word)
}
func (d *dest) popCtxMust(word []byte) {
	l := len(d.ctxStack)
	pop := d.ctxStack[l-1]
	d.ctxStack = d.ctxStack[0 : l-1]
	if pop[0] != word[0] {
		panic(fmt.Errorf("json encoder state: expected to pop %q, popped %q", word, pop))
	}
}
func (d *dest) popCtxMaybe(word []byte) bool {
	l := len(d.ctxStack)
	if l <= 0 {
		return false
	}
	pop := d.ctxStack[l-1]
	if pop[0] != word[0] {
		return false
	}
	d.ctxStack = d.ctxStack[0 : l-1]
	return true
}

// magic words
var (
	wordTrue       = []byte("true")
	wordFalse      = []byte("false")
	wordNull       = []byte("null")
	wordArrayOpen  = []byte("[")
	wordArrayClose = []byte("]")
	wordMapOpen    = []byte("{")
	wordMapClose   = []byte("}")
	wordColon      = []byte(":")
	wordComma      = []byte(",")
)

func (d *dest) OpenMap() {
	d.popCtxMaybe(wordColon)
	if d.popCtxMaybe(wordComma) {
		d.w.Write(wordComma)
	}
	d.pushCtx(wordMapOpen)
	d.w.Write(wordMapOpen)
}
func (d *dest) WriteMapKey(k string) {
	if d.popCtxMaybe(wordComma) {
		d.w.Write(wordComma)
	}
	d.pushCtx(wordColon)
	d.w.Write([]byte(fmt.Sprintf("%q", k)))
	d.w.Write(wordColon)
}
func (d *dest) CloseMap() {
	d.popCtxMaybe(wordComma)
	d.popCtxMust(wordMapOpen)
	d.pushCtx(wordComma)
	d.w.Write(wordMapClose)
}

func (d *dest) OpenArray() {
	d.popCtxMaybe(wordColon)
	if d.popCtxMaybe(wordComma) {
		d.w.Write(wordComma)
	}
	d.pushCtx(wordArrayOpen)
	d.w.Write(wordArrayOpen)
}
func (d *dest) CloseArray() {
	d.popCtxMaybe(wordComma)
	d.popCtxMust(wordArrayOpen)
	d.pushCtx(wordComma)
	d.w.Write(wordArrayClose)
}

func (d *dest) WriteString(v string) {
	d.popCtxMaybe(wordColon)
	if d.popCtxMaybe(wordComma) {
		d.w.Write(wordComma)
	}
	d.pushCtx(wordComma)
	d.w.Write([]byte(fmt.Sprintf("%q", v)))
}

func (d *dest) WriteNull() {
	d.popCtxMaybe(wordColon)
	if d.popCtxMaybe(wordComma) {
		d.w.Write(wordComma)
	}
	d.pushCtx(wordComma)
	d.w.Write(wordNull)
}

// whenever you may be the end of a segment:
/*
	d.popCtxMaybe(wordColon) // in case i'm filling the 'value' slot of a map entry
	if d.popCtxMaybe(wordComma) { // if i'm the non-first entry in an array
		d.w.Write(wordComma) // skip this if you're a map or array open tho -- you're not an end of segment then are you
	}
	d.pushCtx(wordComma) // anyone follow me may need a another comma -- unless you're a map or array end, in which case you'll swallow this
*/
// WriteMapKey is the same, except it can skip the popCtxMaybe(colon).
// TODO : maybe it's weird that we don't branch these maybe-pops on the one-further-up amimaporarray.  we could give better errors by watchdogging that.
