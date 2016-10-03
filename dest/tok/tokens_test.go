package tok

import (
	"fmt"
	"strings"
	"testing"
)

func TestStringing(t *testing.T) {
	assert(t, "token slice fixture should string to expectations",
		strings.Join([]string{
			"Tokens[",
			"\t<{>",
			"\t<k:key>",
			"\t<s:value>",
			"\t<}>",
			"]",
		}, "\n"),
		fmt.Sprintf("%s", Tokens{
			{TokenKind_OpenMap, nil},
			{TokenKind_MapKey, "key"},
			{TokenKind_ValString, "value"},
			{TokenKind_CloseMap, nil},
		}),
	)
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
