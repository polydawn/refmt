package tok

import (
	"fmt"
	"reflect"
)

/*
	Token is either one of magic const tokens (used to denote beginning and
	ending of maps and arrays), or an address to a primitive (string, int, etc).
*/
type Token interface{}

var (
	Token_MapOpen  Token = ctrlToken('{')
	Token_MapClose Token = ctrlToken('}')
	Token_ArrOpen  Token = ctrlToken('[')
	Token_ArrClose Token = ctrlToken(']')
)

/*
	Unexported type used to make sure the control tokens are unique (e.g. so
	you can't accidentally return a rune from buggy code and end up in a very
	strange situation).  Also, attaches a tostring function so you don't
	get a character number in your debug printfs.
*/
type ctrlToken rune

func (t ctrlToken) String() string {
	switch t {
	case Token_MapOpen:
		return "<{>"
	case Token_MapClose:
		return "<}>"
	case Token_ArrOpen:
		return "<[>"
	case Token_ArrClose:
		return "<]>"
	default:
		return "<?>"
	}
}

func IsValidToken(t Token) bool {
	if t == nil {
		return true
	}
	switch t {
	case Token_MapOpen, Token_MapClose, Token_ArrOpen, Token_ArrClose:
		return true
	}
	switch t.(type) {
	case *string, *[]byte:
		return true
	case *bool:
		return true
	case *int, *int8, *int16, *int32, *int64:
		return true
	case *uint, *uint8, *uint16, *uint32, *uint64:
		return true
	}
	return false
}

/*
	Checks if the content of two tokens is the same.
	Tokens are considered the same if they're one of the special consts and are equal;
	or, if they are addresses of a value, then they are the same if they contain the same data
	(it does not matter of the pointers are identical).

	If either value is not a valid token, the result will be false.

	This method is primarily useful for testing.
*/
func IsTokenEqual(t1, t2 Token) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}
	switch t1 {
	case Token_MapOpen, Token_MapClose, Token_ArrOpen, Token_ArrClose:
		return t1 == t2
	}
	switch t2 {
	case Token_MapOpen, Token_MapClose, Token_ArrOpen, Token_ArrClose:
		return false
	}
	if !IsValidToken(t1) {
		return false
	}
	if !IsValidToken(t2) {
		return false
	}
	// Could do another giant type switch here, but don't care about perf much.
	return reflect.ValueOf(t1).Elem().Interface() == reflect.ValueOf(t2).Elem().Interface()
}

func TokenToString(t Token) string {
	switch t {
	case Token_MapOpen:
		return "<{>"
	case Token_MapClose:
		return "<}>"
	case Token_ArrOpen:
		return "<[>"
	case Token_ArrClose:
		return "<]>"
	case nil:
		return "<->"
	}
	if !IsValidToken(t) {
		return fmt.Sprintf("<INVALID:%T:%p>",
			t, t)
	}
	if reflect.ValueOf(t).IsNil() {
		return fmt.Sprintf("<INVALID:%T:%p>",
			t, t)
	}
	return fmt.Sprintf("<%T:%p:%s>",
		t, t, reflect.ValueOf(t).Elem().Interface(),
	)
}

func TokStr(x string) Token { return &x } // Util for testing.
func TokInt(x int) Token    { return &x } // Util for testing.
