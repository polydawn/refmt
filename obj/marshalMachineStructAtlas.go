package obj

import (
	"github.com/polydawn/go-xlate/obj/atlas"
	. "github.com/polydawn/go-xlate/tok"
)

type UnmarshalMachineStructAtlas struct {
	target interface{}
	atlas  atlas.Atlas // Populate on initialization.
	idx    int         // Progress marker
	value  bool        // Progress marker
}

func (m *UnmarshalMachineStructAtlas) Reset(s *Suite, target interface{}) error {
	return nil
}

func (m *UnmarshalMachineStructAtlas) Step(driver *MarshalDriver, s *Suite, tok *Token) (done bool, err error) {
	if m.idx >= len(m.atlas.Fields) {
		panic("incorrect usage: entire struct already walked")
	}
	entry := m.atlas.Fields[m.idx]
	valp := entry.Grab(m.target)
	if m.value {
		driver.Recurse(tok, valp, s.pickMarshalMachine(valp))
	} else {
		*tok = entry.Name
	}
	m.value = !m.value
	return true, nil
}
