package obj

import (
	"fmt"
	"reflect"

	"github.com/polydawn/go-xlate/obj/atlas"
	. "github.com/polydawn/go-xlate/tok"
)

type MarshalMachineStructAtlas struct {
	target interface{}
	atlas  atlas.Atlas // Populate on initialization.
	index  int         // Progress marker
	value  bool        // Progress marker
}

func NewMarshalMachineStructAtlas(atl atlas.Atlas) MarshalMachine {
	atl.Init()
	return &MarshalMachineStructAtlas{atlas: atl}
}

func (m *MarshalMachineStructAtlas) Reset(s *Suite, target interface{}) error {
	m.target = target
	m.index = -1
	m.value = false
	if !reflect.ValueOf(target).CanAddr() {
		return fmt.Errorf("error resetting MarshalMachineStructAtlas: target is not addressable")
	}
	return nil
}

func (m *MarshalMachineStructAtlas) Step(driver *MarshalDriver, s *Suite, tok *Token) (done bool, err error) {
	if m.index < 0 {
		if m.target == nil { // REVIEW p sure should have ptr cast and indirect
			*tok = nil
			m.index++
			return true, nil
		}
		*tok = Token_MapOpen
		m.index++
		return false, nil
	}
	nEntries := len(m.atlas.Fields)
	if m.index == nEntries {
		*tok = Token_MapClose
		m.index++
		return true, nil
	}
	if m.index > nEntries {
		return true, fmt.Errorf("invalid state: entire struct (%d fields) already consumed", nEntries)
	}

	entry := m.atlas.Fields[m.index]
	if m.value {
		valp := entry.Grab(m.target)
		m.index++
		m.value = false
		return false, driver.Recurse(tok, valp, s.pickMarshalMachine(valp))
	}
	*tok = &entry.Name
	m.value = true
	return false, nil
}
