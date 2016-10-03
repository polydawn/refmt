package xlate

import (
	"fmt"
	"reflect"
)

/*
	Maps the input to a string as per printf'ing with `"%s"`.
*/
func Map_Wildcard_toString(_ *Mapper, dest Destination, input interface{}) {
	dest.WriteString(fmt.Sprintf("%s", input))
}

var _ MappingFunc = Map_Wildcard_toString

/*
	Maps the input to a string of simply the type name.
*/
func Map_Wildcard_toStringOfType(_ *Mapper, dest Destination, input interface{}) {
	dest.WriteString(reflect.TypeOf(input).Name())
}

var _ MappingFunc = Map_Wildcard_toStringOfType
