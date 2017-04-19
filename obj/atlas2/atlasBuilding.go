package atlas

import "reflect"

func Build(entries ...AtlasEntry) (Atlas, error) {
	return Atlas{}, nil // TODO
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
