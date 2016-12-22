package obj

import (
	"fmt"

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
	return &MarshalMachineStructAtlas{atlas: atl}
}

func (m *MarshalMachineStructAtlas) Reset(s *Suite, target interface{}) error {
	m.target = target
	m.index = -1
	m.value = false
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
	if m.index == len(m.atlas.Fields) {
		*tok = Token_MapClose
		m.index++
		return true, nil
	}
	if m.index > len(m.atlas.Fields) {
		return true, fmt.Errorf("invalid state: entire struct already consumed")
	}

	entry := m.atlas.Fields[m.index]
	if m.value {
		//fmt.Printf(">> %d : %#T %#v\n   - : %#T %#v\n", m.index, m.target, m.target, *(m.target.(*interface{})), *(m.target.(*interface{})))
		valp := entry.Grab(m.target)
		//fmt.Printf(":: %d : %#T %#v\n   - : %#T %#v\n", m.index, valp, valp, *(valp.(*interface{})), *(valp.(*interface{})))
		m.index++
		return false, driver.Recurse(tok, valp, s.pickMarshalMachine(valp))
	}
	*tok = &entry.Name
	m.value = true
	return false, nil
}
