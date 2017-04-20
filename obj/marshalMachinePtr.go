package obj

import (
	"reflect"

	. "github.com/polydawn/refmt/tok"
)

type ptrDerefDelegateMarshalMachine struct {
	MarshalMachine
	peelCount int

	isNil bool
}

func (m *ptrDerefDelegateMarshalMachine) Reset(s *marshalSlab, valp interface{}) error {
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

func (m *ptrDerefDelegateMarshalMachine) Step(d *MarshalDriver, s *marshalSlab, tok *Token) (done bool, err error) {
	if m.isNil {
		tok.Type = TNull
		return true, nil
	}
	return m.MarshalMachine.Step(d, s, tok)
}
