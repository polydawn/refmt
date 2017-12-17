package obj

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/polydawn/refmt/tok"
)

type marshalMachineMapWildcard struct {
	cfg *atlas.AtlasEntry // set on initialization

	target_rv reflect.Value
	value_rt  reflect.Type
	valueMach MarshalMachine
	keys      []wildcardMapStringyKey
	index     int
	value     bool
}

func (mach *marshalMachineMapWildcard) Reset(slab *marshalSlab, rv reflect.Value, rt reflect.Type) error {
	mach.target_rv = rv

	// Pick machinery for handling the value types.
	mach.value_rt = rt.Elem()
	mach.valueMach = slab.requisitionMachine(mach.value_rt)

	// Enumerate all the keys (must do this up front, one way or another),
	// flip them into strings,
	// and sort them (optional, arguably, but right now you're getting it).
	key_rt := rt.Key()
	switch key_rt.Kind() {
	case reflect.String:
		// continue.
		// note: stdlib json.marshal supports all the int types here as well, and will
		//  tostring themach.  but this is not supported symmetrically; so we simply... don't.
		//  we could also consider supporting anything that uses a MarshalTransformFunc
		//  to become a string kind; that's a fair bit of code, perhaps later.
	default:
		return fmt.Errorf("unsupported map key type %q", key_rt.Name())
	}
	keys_rv := mach.target_rv.MapKeys()
	mach.keys = make([]wildcardMapStringyKey, len(keys_rv))
	for i, v := range keys_rv {
		mach.keys[i].rv = v
		mach.keys[i].s = v.String()
	}
	switch mach.cfg.MapMorphism.KeySortMode {
	case atlas.KeySortMode_Default:
		sort.Sort(wildcardMapStringyKey_byString(mach.keys))
	case atlas.KeySortMode_RFC7049:
		sort.Sort(wildcardMapStringyKey_RFC7049(mach.keys))
	default:
		panic(fmt.Errorf("unknown map key sort mode %q", mach.cfg.MapMorphism.KeySortMode))
	}

	mach.index = -1
	return nil
}

func (mach *marshalMachineMapWildcard) Step(driver *Marshaller, slab *marshalSlab, tok *Token) (done bool, err error) {
	if mach.index < 0 {
		if mach.target_rv.IsNil() {
			tok.Type = TNull
			mach.index++
			return true, nil
		}
		tok.Type = TMapOpen
		tok.Length = mach.target_rv.Len()
		mach.index++
		return false, nil
	}
	if mach.index == len(mach.keys) {
		tok.Type = TMapClose
		mach.index++
		slab.release()
		return true, nil
	}
	if mach.index > len(mach.keys) {
		return true, fmt.Errorf("invalid state: value already consumed")
	}
	if mach.value {
		val_rv := mach.target_rv.MapIndex(mach.keys[mach.index].rv)
		mach.value = false
		mach.index++
		return false, driver.Recurse(tok, val_rv, mach.value_rt, mach.valueMach)
	}
	tok.Type = TString
	tok.Str = mach.keys[mach.index].s
	mach.value = true
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

type wildcardMapStringyKey_RFC7049 []wildcardMapStringyKey

func (x wildcardMapStringyKey_RFC7049) Len() int      { return len(x) }
func (x wildcardMapStringyKey_RFC7049) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x wildcardMapStringyKey_RFC7049) Less(i, j int) bool {
	li, lj := len(x[i].s), len(x[j].s)
	if li == lj {
		return x[i].s < x[j].s
	}
	return li < lj
}
