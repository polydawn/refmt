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
	Token_MapOpen  Token = '{'
	Token_MapClose Token = '}'
	Token_ArrOpen  Token = '['
	Token_ArrClose Token = ']'
)

func IsValidToken(t Token) bool {
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
	switch t1 {
	case Token_MapOpen, Token_MapClose, Token_ArrOpen, Token_ArrClose:
		return t1 == t2
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
	}
	if !IsValidToken(t) {
		return "<???>"
	}
	return fmt.Sprintf("<%T:%p:%s>",
		t, t, reflect.ValueOf(t).Elem().Interface(),
	)
}
