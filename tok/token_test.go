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
		{'[', true},  // ...!  Remember, casting makes it so.
		{"{", false}, // ... but strings aren't so casted.

		// Some actual values:
		{&str, true},
		{"", false},
		{&i, true},
		{4, false},
	}
	for _, tr := range tt {
		Assert(t, fmt.Sprintf("validity check for %#v", tr.tok),
			tr.valid, IsValidToken(tr.tok))
	}
}

func TestTokenEqualityDefn(t *testing.T) {
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
	}
	for _, tr := range tt {
		Assert(t, fmt.Sprintf("equality check for %#v==%#v", tr.tok1, tr.tok2),
			tr.eq, IsTokenEqual(tr.tok1, tr.tok2))
	}
}
