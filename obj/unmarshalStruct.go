package obj

import (
	"fmt"
	"reflect"

	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/polydawn/refmt/tok"
)

type unmarshalMachineStructAtlas struct {
	cfg *atlas.AtlasEntry // set on initialization

	rv         reflect.Value
	expectLen  int                  // Length header from mapOpen token.  If it was set, we validate it.
	index      int                  // Progress marker: our distance into the stream of pairs.
	value      bool                 // Progress marker: whether the next token is a value.
	fieldEntry atlas.StructMapEntry // Which field we expect next: set when consuming a key.
}

func (mach *unmarshalMachineStructAtlas) Reset(_ *unmarshalSlab, rv reflect.Value, _ reflect.Type) error {
	mach.rv = rv
	// not necessary to reset expectLen because MapOpen tokens also consistently use the -1 convention.
	mach.index = -1
	mach.value = false
	return nil
}

func (mach *unmarshalMachineStructAtlas) Step(driver *Unmarshaller, slab *unmarshalSlab, tok *Token) (done bool, err error) {
	// Starter state.
	if mach.index < 0 {
		switch tok.Type {
		case TMapOpen:
			// Great.  Consumed.
			mach.expectLen = tok.Length
			mach.index++
			return false, nil
		case TMapClose:
			return true, ErrMalformedTokenStream{tok.Type, "start of map"}
		case TArrOpen:
			return true, ErrMalformedTokenStream{tok.Type, "start of map"}
		case TArrClose:
			return true, ErrMalformedTokenStream{tok.Type, "start of map"}
		case TNull:
			mach.rv.Set(reflect.Zero(mach.rv.Type()))
			return true, nil
		default:
			return true, ErrMalformedTokenStream{tok.Type, "start of map"}
		}
	}

	// Accept value:
	if mach.value {
		child_rv := mach.fieldEntry.ReflectRoute.TraverseToValue(mach.rv)
		mach.index++
		mach.value = false
		return false, driver.Recurse(
			tok,
			child_rv,
			mach.fieldEntry.Type,
			slab.requisitionMachine(mach.fieldEntry.Type),
		)
	}

	// Accept key or end:
	if mach.index > 0 {
		slab.release()
	}
	switch tok.Type {
	case TMapClose:
		// If we got length header, validate that; error if mismatch.
		if mach.expectLen >= 0 {
			if mach.expectLen != mach.index {
				return true, fmt.Errorf("malformed map token stream: declared length %d, actually got %d entries", mach.expectLen, mach.index)
			}
		}

		// Future: this would be a reasonable place to check that all required fields have been filled in, if we add such a feature.

		return true, nil
	case TString:
		for n := 0; n < len(mach.cfg.StructMap.Fields); n++ {
			fieldEntry := mach.cfg.StructMap.Fields[n]
			if fieldEntry.SerialName != tok.Str {
				continue
			}
			mach.fieldEntry = fieldEntry
			mach.value = true
			break
		}
		if mach.value == false {
			// FUTURE: it should be configurable per atlas.StructMap whether this is considered an error or to be tolerated.
			// Currently we're being extremely strict about it, which is a divergence from the stdlib json behavior.
			return true, ErrNoSuchField{tok.Str, mach.cfg.Type.String()}
		}
	default:
		return true, ErrMalformedTokenStream{tok.Type, "map key"}
	}
	return false, nil
}
