package testutil

// Local drop-in for github.com/warpfork/go-wish's Wish/ShouldEqual API.
// Equivalent in shape and behaviour, but uses reflect.DeepEqual instead of
// go-wish's vendored cmp package. The cmp package does unsafe pointer
// arithmetic in retrieveUnexportedField that Go's checkptr (enabled with
// -race) flags as invalid on Go 1.25+. go-wish is unmaintained, so this
// shim lets refmt's tests keep running under -race without changing call
// sites. Restore the upstream import if go-wish is ever fixed.

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type Checker func(actual, desired interface{}) (problem string, passed bool)

func Wish(t *testing.T, actual interface{}, check Checker, desired interface{}) bool {
	t.Helper()
	problem, passed := check(actual, desired)
	if !passed {
		t.Logf("check rejected:\n%s", indent(problem))
		t.Fail()
	}
	return passed
}

func Require(t *testing.T, actual interface{}, check Checker, desired interface{}) {
	t.Helper()
	problem, passed := check(actual, desired)
	if !passed {
		t.Logf("halting: critical check rejected:\n%s", indent(problem))
		t.FailNow()
	}
}

func ShouldEqual(actual, desired interface{}) (string, bool) {
	if reflect.DeepEqual(actual, desired) {
		return "", true
	}
	return fmt.Sprintf("expected: %#v\nactual:   %#v", desired, actual), false
}

func indent(s string) string {
	if s == "" {
		return "\t"
	}
	return "\t" + strings.ReplaceAll(strings.TrimRight(s, "\n"), "\n", "\n\t")
}
