package xlate

import (
	"fmt"
)

/*
	Error raised by `NewMapper` if two or more rows in the provided
	`MapperSetup` are for the same types.
*/
type ErrMapperSetupNotUnique struct {
	Row MapperSetupRow
}

func (e ErrMapperSetupNotUnique) Error() string {
	return fmt.Sprintf("ErrMapperSetupNotUnique {row:%#v}", e.Row)
}

/*
	Error raised when expecting a MappingFunc and finding a nil.
*/
type ErrNilMappingFunc struct {
	Row MapperSetupRow
}

func (e ErrNilMappingFunc) Error() string {
	return fmt.Sprintf("ErrNilMappingFunc {row:%#v}", e.Row)
}
