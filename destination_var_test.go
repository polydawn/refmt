package xlate

import (
	"testing"
)

func TestDestinationVar(t *testing.T) {
	var str string
	d := NewVarDestination(&str)
	d.WriteString("test")

	assert(t, "output to string var, input a string", "test", str)
}
