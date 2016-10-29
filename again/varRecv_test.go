package again

import (
	"fmt"
	"testing"
)

func TestWow(t *testing.T) {
	var v int
	vr := NewVarReceiver(&v)
	stepShouldDone(t, vr, 4)
	assert(t, "simple literal test", 4, v)

	var v2 interface{}
	vr = NewVarReceiver(&v2)
	stepShouldContinue(t, vr, Token_MapOpen)
	stepShouldDone(t, vr, Token_MapClose)
	assert(t, "map and recurse test",
		map[string]interface{}{},
		v2)

	var v3 interface{}
	vr = NewVarReceiver(&v3)
	stepShouldContinue(t, vr, Token_MapOpen)
	stepShouldContinue(t, vr, "key")
	stepShouldContinue(t, vr, 6)
	stepShouldDone(t, vr, Token_MapClose)
	assert(t, "map and recurse test",
		map[string]interface{}{"key": 6},
		v3)
}

func stringyEquality(x, y interface{}) bool {
	return fmt.Sprintf("%#v", x) == fmt.Sprintf("%#v", y)
}

func assert(t *testing.T, title string, expect, actual interface{}) {
	if !stringyEquality(expect, actual) {
		t.Errorf("test %q FAILED:\n\texpected  %#v\n\tactual    %#v",
			title, expect, actual)
	}
}

func stepShouldContinue(t *testing.T, sink TokenSink, tok Token) {
	done, err := sink.Step(&tok)
	if err != nil {
		t.Errorf("step errored: %s", err)
	}
	if done {
		t.Errorf("expected step not to be done")
	}
}

func stepShouldDone(t *testing.T, sink TokenSink, tok Token) {
	done, err := sink.Step(&tok)
	if err != nil {
		t.Errorf("step errored: %s", err)
	}
	if !done {
		t.Errorf("expected step be done")
	}
}
