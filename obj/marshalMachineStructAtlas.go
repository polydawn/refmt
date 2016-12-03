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

func (m *UnmarshalMachineStructAtlas) init() {

}

func (m *UnmarshalMachineStructAtlas) Step(driver *MarshalDriver, tok *Token) (done bool, err error) {
	if m.idx >= len(m.atlas.Fields) {
		panic("incorrect usage: entire struct already walked")
	}
	entry := m.atlas.Fields[m.idx]
	if m.value {
		driver.Recurse(tok, entry.Grab(m.target))
	} else {
		*tok = entry.Name
	}
	m.value = !m.value
	return true, nil
}
