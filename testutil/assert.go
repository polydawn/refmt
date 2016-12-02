package testutil

import (
	"fmt"
	"testing"
)

func stringyEquality(x, y interface{}) bool {
	return fmt.Sprintf("%#v", x) == fmt.Sprintf("%#v", y)
}

func Assert(t *testing.T, title string, expect, actual interface{}) {
	if !stringyEquality(expect, actual) {
		t.Errorf("FAILED test %q:\n\texpected  %#v\n\tactual    %#v",
			title, expect, actual)
	} else {
		t.Logf("test %q result matched", title)
	}
}

func capturePanics(fn func()) (e error) {
	defer func() {
		if rcvr := recover(); rcvr != nil {
			e = rcvr.(error)
		}
	}()
	fn()
	return
}
