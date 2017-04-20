package atlas

import (
	"fmt"
	"reflect"
)

func Build(entries ...AtlasEntry) (Atlas, error) {
	atl := Atlas{
		mappings: make(map[uintptr]*AtlasEntry),
	}
	for _, entry := range entries {
		rtid := reflect.ValueOf(entry.Type).Pointer()
		if _, exists := atl.mappings[rtid]; exists {
			return Atlas{}, fmt.Errorf("repeated entry for %v", entry.Type)
		}
		atl.mappings[rtid] = &entry
	}
	return atl, nil
}

func BuildEntry(typeHintObj interface{}) *BuilderCore {
	return &BuilderCore{
		&AtlasEntry{Type: reflect.TypeOf(typeHintObj)},
	}
}

/*
	Intermediate step in building an AtlasEntry: use `BuildEntry` to
	get one of these to start with, then call one of the methods
	on this type to get a specialized builder which has the methods
	relevant for setting up that specific kind of mapping.
*/
type BuilderCore struct {
	entry *AtlasEntry
}
