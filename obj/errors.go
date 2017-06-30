package obj

import (
	"fmt"
	"reflect"

	. "github.com/polydawn/refmt/tok"
)

// ErrInvalidUnmarshalTarget describes an invalid argument passed to UnmarshalDriver.Bind.
// (Unmarshalling must target a non-nil pointer so that it can address the value.)
type ErrInvalidUnmarshalTarget struct {
	Type reflect.Type
}

func (e ErrInvalidUnmarshalTarget) Error() string {
	if e.Type == nil {
		return "unmarshal error: invalid target (nil)"
	}
	if e.Type.Kind() != reflect.Ptr {
		return "unmarshal error: invalid target (non-pointer " + e.Type.String() + ")"
	}
	return "unmarshal error: invalid target (nil " + e.Type.String() + ")"
}

// ErrUnmarshalTypeCantFit is the error returned when unmarshalling cannot
// coerce the tokens in the stream into the kind of variables the unmarshal is targetting,
// for example if a map open token comes when an int is expected,
// or an int token comes when a string is expected.
type ErrUnmarshalTypeCantFit struct {
	Token Token
	Value reflect.Value
}

func (e ErrUnmarshalTypeCantFit) Error() string {
	return fmt.Sprintf("unmarshal error: cannot assign %s to %s field", e.Token, e.Value.Kind())
}

// ErrMalformedTokenStream is the error returned when unmarshalling recieves a
// completely invalid transition, such as when a map value is expected, but the
// map suddenly closes, or an array close is recieved with no matching array open.
type ErrMalformedTokenStream struct {
	Got      TokenType // Token in the stream that triggered the error.
	Expected string    // Freeform string describing valid token types.  Often a summary like "array close or start of value", or "map close or key".
}

func (e ErrMalformedTokenStream) Error() string {
	return fmt.Sprintf("malformed stream: invalid appearance of %s token; expected %s", e.Got, e.Expected)
}

// ErrNoSuchField is the error returned when unmarshalling into a struct and
// the token stream for the map contains a key which is not defined for the struct.
type ErrNoSuchField struct {
	Name string // Field name from the token.
}

func (e ErrNoSuchField) Error() string {
	return fmt.Sprintf("unmarshal error: no such field named %s", e.Name)
}
