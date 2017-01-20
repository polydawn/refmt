package obj

import (
	"fmt"
	"reflect"
	"sort"

	. "github.com/polydawn/go-xlate/tok"
)

type MarshalMachineMapWildcard struct {
	target_rv reflect.Value
	valueMach MarshalMachine
	keys      []wildcardMapStringyKey
	index     int
	value     bool
}

func (m *MarshalMachineMapWildcard) Reset(s *Suite, valp interface{}) error {
	m.target_rv = reflect.ValueOf(valp).Elem()

	// Pick machinery for handling the value types.
	m.valueMach = s.marshalMachineForType(m.target_rv.Type().Elem())

	// Enumerate all the keys (must do this up front, one way or another),
	// flip them into strings,
	// and sort them (optional, arguably, but right now you're getting it).
	key_rt := m.target_rv.Type().Key()
	switch key_rt.Kind() {
	case reflect.String:
		// continue.
		// note: stdlib json.marshal supports all the int types here as well, and will
		//  tostring them.  but this is not supported symmetrically; so we simply... don't.
	default:
		return fmt.Errorf("unsupported map key type %q", key_rt.Name())
	}
	keys_rv := m.target_rv.MapKeys()
	m.keys = make([]wildcardMapStringyKey, len(keys_rv))
	for i, v := range keys_rv {
		m.keys[i].rv = v
		m.keys[i].s = v.String()
	}
	sort.Sort(wildcardMapStringyKey_byString(m.keys))

	m.index = -1
	return nil
}

func (m *MarshalMachineMapWildcard) Step(d *MarshalDriver, s *Suite, tok *Token) (done bool, err error) {
	if m.index < 0 {
		if m.target_rv.IsNil() {
			*tok = nil
			m.index++
			return true, nil
		}
		*tok = Token_MapOpen
		m.index++
		return false, nil
	}
	if m.index == len(m.keys) {
		*tok = Token_MapClose
		m.index++
		return true, nil
	}
	if m.index > len(m.keys) {
		return true, fmt.Errorf("invalid state: value already consumed")
	}
	if m.value {
		val_rv := m.target_rv.MapIndex(m.keys[m.index].rv)
		new_vprv := reflect.New(val_rv.Type())
		new_vprv.Elem().Set(val_rv)
		valp := new_vprv.Interface()
		m.value = false
		m.index++
		return false, d.Recurse(tok, valp, s.pickMarshalMachine(valp))
	}
	*tok = &(m.keys[m.index].s)
	m.value = true
	return false, nil
}

// Holder for the reflect.Value and string form of a key.
// We need the reflect.Value for looking up the map value;
// and we need the string for sorting.
type wildcardMapStringyKey struct {
	rv reflect.Value
	s  string
}

type wildcardMapStringyKey_byString []wildcardMapStringyKey

func (x wildcardMapStringyKey_byString) Len() int           { return len(x) }
func (x wildcardMapStringyKey_byString) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x wildcardMapStringyKey_byString) Less(i, j int) bool { return x[i].s < x[j].s }
