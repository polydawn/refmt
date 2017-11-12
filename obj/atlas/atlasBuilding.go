package atlas

import (
	"fmt"
	"reflect"
)

func Build(entries ...*AtlasEntry) (Atlas, error) {
	atl := Atlas{
		mappings:    make(map[uintptr]*AtlasEntry),
		tagMappings: make(map[int]*AtlasEntry),
	}
	for _, entry := range entries {
		rtid := reflect.ValueOf(entry.Type).Pointer()
		if _, exists := atl.mappings[rtid]; exists {
			return Atlas{}, fmt.Errorf("repeated entry for type %v", entry.Type)
		}
		atl.mappings[rtid] = entry

		if entry.Tagged == true {
			if prev, exists := atl.tagMappings[entry.Tag]; exists {
				return Atlas{}, fmt.Errorf("repeated tag %v on type %v (already mapped to type %v)", entry.Tag, entry.Type, prev.Type)
			}
			atl.tagMappings[entry.Tag] = entry
		}
	}
	return atl, nil
}
func MustBuild(entries ...*AtlasEntry) Atlas {
	atl, err := Build(entries...)
	if err != nil {
		panic(err)
	}
	return atl
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

	One full example of using this builder may look like the following:

		atlas.BuildEntry(Formula{}).StructMap().Autogenerate().Complete()

	Some intermediate manipulations may be performed on this object,
	for example setting the "tag" (if you want to use cbor tagging),
	before calling the specializer method.
	In this case, just keep chaining the configuration calls like so:

		atlas.BuildEntry(Formula{}).UseTag(4000)
			.StructMap().Autogenerate().Complete()

*/
type BuilderCore struct {
	entry *AtlasEntry
}

func (x *BuilderCore) UseTag(tag int) *BuilderCore {
	x.entry.Tagged = true
	x.entry.Tag = tag
	return x
}
