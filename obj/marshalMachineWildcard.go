package obj

import (
	"reflect"

	. "github.com/polydawn/refmt/tok"
)

/*
	A MarshalMachine unwraps an `interface{}` value,
	selects the correct machinery for handling its content,
	and delegates immediately to that machine.
*/
type MarshalMachineWildcard struct {
	delegate MarshalMachine
}

func (m *MarshalMachineWildcard) Reset(s *marshalSlab, valp interface{}) error {
	val_rv := reflect.ValueOf(valp).Elem()
	// If the interface contains nil, go no further; we'll simply yield that single token.
	if val_rv.IsNil() {
		return nil
	}
	// Pick, reset, and retain a delegate machine for the interior type.
	val_unwrap_rv := val_rv.Elem() // unwrap iface
	m.delegate = s.mustPickMarshalMachineByType(val_unwrap_rv.Type())
	// Values stored in an interface wildcard are not themselves addressable, so, some jiggerypokery required:
	new_vprv := reflect.New(val_unwrap_rv.Type())
	new_vprv.Elem().Set(val_unwrap_rv)
	valp = new_vprv.Interface()
	return m.delegate.Reset(s, valp)
}

func (m MarshalMachineWildcard) Step(driver *MarshalDriver, s *marshalSlab, tok *Token) (done bool, err error) {
	if m.delegate == nil {
		tok.Type = TNull
		return true, nil
	}
	return m.delegate.Step(driver, s, tok)
}
