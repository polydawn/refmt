package obj

import (
	"fmt"

	"github.com/polydawn/go-xlate/obj/atlas"
	. "github.com/polydawn/go-xlate/tok"
)

type MarshalMachineStructAtlas struct {
	atlas atlas.Atlas // Populate on initialization.

	target interface{}
	index  int  // Progress marker
	value  bool // Progress marker
}

func (m *MarshalMachineStructAtlas) Reset(s *slab, target interface{}) error {
	m.target = target
	m.index = -1
	m.value = false
	return nil
}

func (m *MarshalMachineStructAtlas) Step(driver *MarshalDriver, s *slab, tok *Token) (done bool, err error) {
	//fmt.Printf("--step on %#v: i=%d/%d v=%v\n", m.target, m.index, len(m.atlas.Fields), m.value)
	if m.index < 0 {
		if m.target == nil {
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
		s.release()
		return true, nil
	}
	if m.index > nEntries {
		return true, fmt.Errorf("invalid state: entire struct (%d fields) already consumed", nEntries)
	}

	if m.value {
		valp := m.atlas.Fields[m.index].Grab(m.target)
		m.index++
		m.value = false
		return false, driver.Recurse(tok, valp, s.mustPickMarshalMachine(valp))
	}
	*tok = &m.atlas.Fields[m.index].Name
	m.value = true
	if m.index > 0 {
		s.release()
	}
	return false, nil
}
