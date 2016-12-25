package tok

import (
	"fmt"
	"testing"

	. "github.com/polydawn/go-xlate/testutil"
)

func TestTokenValidityDefn(t *testing.T) {
	var str string
	var i int
	tt := []struct {
		tok   Token
		valid bool
	}{
		// The control constants:
		{Token_MapOpen, true},
		{Token_MapClose, true},
		{Token_ArrOpen, true},
		{Token_ArrClose, true},
		{'[', false}, // casting will not make it so, due to the unexported typedef.
		{"{", false},

		// Some actual values:
		{&str, true},
		{"", false},
		{&i, true},
		{4, false},
	}
	for _, tr := range tt {
		Assert(t, fmt.Sprintf("validity check for %q", tr.tok),
			tr.valid, IsValidToken(tr.tok))
	}
}

func TestTokenEqualityDefn(t *testing.T) {
	var str1 string = "one"
	var str2 string = "two"
	var str1b string = "one"
	var i3 int = 3
	var i4 int = 4
	var i3b int = 3
	tt := []struct {
		tok1 Token
		tok2 Token
		eq   bool
	}{
		// The control constants must equal themselves(!):
		{Token_MapOpen, Token_MapOpen, true},
		{Token_MapClose, Token_MapClose, true},
		{Token_ArrOpen, Token_ArrOpen, true},
		{Token_ArrClose, Token_ArrClose, true},

		// The control constants must not equal each other:
		{Token_MapOpen, Token_MapClose, false},
		{Token_MapOpen, Token_ArrOpen, false},
		{Token_ArrOpen, Token_ArrClose, false},
		{Token_ArrOpen, Token_MapClose, false},

		// Other values should behave as you expect:
		{&str1, &str1, true},  // self is same
		{&str1, &str2, false}, // other content is not
		{&str1, &str1b, true}, // other ptr, but same content is same
		{&i3, &i3, true},      // self is same
		{&i3, &i4, false},     // other content is not
		{&i3, &i3b, true},     // other ptr, but same content is same
		{&i3, &str1, false},   // totally different types are not same
		{&i3, Token_MapOpen, false},
		{Token_MapOpen, &i3, false},
	}
	for _, tr := range tt {
		Assert(t, fmt.Sprintf("equality check for %q==%q", tr.tok1, tr.tok2),
			tr.eq, IsTokenEqual(tr.tok1, tr.tok2))
	}
}
