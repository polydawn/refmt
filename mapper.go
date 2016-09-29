package xlate

import (
	"reflect"
)

/*
	A MappingFunc accepts an `input` object, and should put its contents
	into the `dest` interface by calling the relevant methods to place each
	field as a map, or each row in an array.
*/
type MappingFunc func(dest Destination, input interface{})

/*


	(Implementation note: this is a slice of config hunks rather than a map,
	because there are limitations on the types of map keys; we don't need to
	be bound by most of those limitations.)
*/
type MapperSetup []MapperSetupRow

type MapperSetupRow struct {
	Type interface{}
	Func MappingFunc
}

type Mapper struct {
	mappings map[reflect.Type]MappingFunc
}

/*
	Initialize a new Mapper.

	NewMapper is used with inputs that are pretty much impossible to imagine
	being user input, and so any errors are panics: they shouldn't have been
	here at program compile time.

	May panic with:

	  - `*xlate.ErrMapperSetupNotUnique`
	  - `*xlate.ErrNilMappingFunc`
	  - anything terrible with reflection
*/
func NewMapper(mappings MapperSetup) *Mapper {
	return &Mapper{
		mappings: processMapperSetup(mappings),
	}
}
func processMapperSetup(ms MapperSetup) map[reflect.Type]MappingFunc {
	mappings := make(map[reflect.Type]MappingFunc, len(ms))
	for _, row := range ms {
		// Translate reflection.
		typ := reflect.TypeOf(row.Type)
		// Validate.
		if row.Func == nil {
			panic(&ErrNilMappingFunc{row})
		}
		if mappings[typ] != nil {
			panic(&ErrMapperSetupNotUnique{row})
		}
		// Accept.
		mappings[typ] = row.Func
	}
	return mappings
}

func (m *Mapper) Map(dest Destination, input interface{}) {
	// TODO run down the setup table for matches
}
