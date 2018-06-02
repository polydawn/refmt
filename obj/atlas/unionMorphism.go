package atlas

import (
	"fmt"
	"reflect"
)

type UnionKeyedMorphism struct {
	// Mapping of typehint key strings to atlasEntry that should be delegated to.
	Elements map[string]*AtlasEntry
	// Purely to have in readiness for error messaging.
	KnownMembers []string
}

func (x *BuilderCore) KeyedUnion() *BuilderUnionKeyedMorphism {
	if x.entry.Type.Kind() != reflect.Interface {
		panic(fmt.Errorf("cannot use union morphisms for type %q, which is kind %s", x.entry.Type, x.entry.Type.Kind()))
	}
	x.entry.UnionKeyedMorphism = &UnionKeyedMorphism{}
	return &BuilderUnionKeyedMorphism{x.entry}
}

type BuilderUnionKeyedMorphism struct {
	entry *AtlasEntry
}

func (x *BuilderUnionKeyedMorphism) Of(elements map[string]*AtlasEntry) *AtlasEntry {
	x.entry.UnionKeyedMorphism.Elements = elements
	// FIXME: do a copy loop plus sanity check that all the delegates are... well struct or map machines really, but definitely blacklisting other delegating machinery.
	// FIXME: and sanity check that they can all be assigned to the interface ffs.
	// FIXME: populate KnownMembers at the same time.
	return x.entry
}
