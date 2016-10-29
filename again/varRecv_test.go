package again

import (
	"fmt"
	"testing"
)

func TestWow(t *testing.T) {
	var v int
	vr := NewVarReceiver(&v)
	tok := Token(4)
	stepShouldDone(t, vr, &tok)
	assert(t, "simple literal test", 4, v)

	var v2 interface{}
	vr = NewVarReceiver(&v2)
	stepShouldContinue(t, vr, &Token_MapOpen)
	stepShouldDone(t, vr, &Token_MapClose)
	assert(t, "map and recurse test",
		map[string]interface{}{},
		v2)

	var v3 interface{}
	vr = NewVarReceiver(&v3)
	stepShouldContinue(t, vr, &Token_MapOpen)
	tok = "key"
	stepShouldContinue(t, vr, &tok)
	tok = 6
	stepShouldContinue(t, vr, &tok)
	stepShouldDone(t, vr, &Token_MapClose)
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

func stepShouldContinue(t *testing.T, sink TokenSink, tok *Token) {
	done, err := sink.Step(tok)
	if err != nil {
		t.Errorf("step errored: %s", err)
	}
	if done {
		t.Errorf("expected step not to be done")
	}
}

func stepShouldDone(t *testing.T, sink TokenSink, tok *Token) {
	done, err := sink.Step(tok)
	if err != nil {
		t.Errorf("step errored: %s", err)
	}
	if !done {
		t.Errorf("expected step be done")
	}
}
