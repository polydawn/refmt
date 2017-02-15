package obj

import (
	"reflect"

	"github.com/polydawn/go-xlate/tok"
)

type ptrDerefDelegateMarshalMachine struct {
	MarshalMachine
	peelCount int

	isNil bool
}

func (m *ptrDerefDelegateMarshalMachine) Reset(s *slab, valp interface{}) error {
	m.isNil = false
	rv := reflect.ValueOf(valp)
	for i := 0; i < m.peelCount; i++ {
		rv = rv.Elem()
		if rv.IsNil() {
			m.isNil = true
			return nil
		}
	}
	return m.MarshalMachine.Reset(s, rv.Interface())
}

func (m *ptrDerefDelegateMarshalMachine) Step(d *MarshalDriver, s *slab, tok *tok.Token) (done bool, err error) {
	if m.isNil {
		*tok = nil
		return true, nil
	}
	return m.MarshalMachine.Step(d, s, tok)
}
