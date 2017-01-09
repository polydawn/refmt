package obj

import (
	"reflect"

	"github.com/polydawn/go-xlate/tok"
)

/*
	A MarshalMachine which handles (any) pointer indirection.
*/
type MarshalMachinePtr struct {
	target interface{}
}

func (m *MarshalMachinePtr) Reset(_ *Suite, target interface{}) error {
	m.target = target
	return nil
}

func (m MarshalMachinePtr) Step(d *MarshalDriver, s *Suite, tok *tok.Token) (done bool, err error) {
	//	fmt.Printf(":hai -- ptr fondling %#v\n", m.target)
	val_rv := reflect.ValueOf(m.target).Elem()
	if val_rv.IsNil() {
		*tok = nil
		return true, nil
	}
	derefvalp_rv := val_rv.Elem().Addr()
	//	fmt.Printf(":: unwound to %#v\n", derefvalp_rv.Interface())
	return true, d.Recurse(tok, derefvalp_rv.Interface(), s.pickMarshalMachine(derefvalp_rv.Interface()))
}
