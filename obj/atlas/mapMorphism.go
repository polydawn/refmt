package atlas

import (
	"fmt"
)

// A type to enumerate key sorting modes.
type KeySortMode string

const (
	KeySortMode_Default = "default" // e.g. lexical string sort for strings, etc
	KeySortMode_RFC7049 = "rfc7049" // "Canonical" as proposed by rfc7049 § 3.9 (shorter byte sequences sort to top)
)

type MapMorphism struct {
	KeySortMode KeySortMode
}

func (x *BuilderCore) MapMorphism() *BuilderMapMorphism {
	x.entry.MapMorphism = &MapMorphism{
		KeySortMode_Default,
	}
	return &BuilderMapMorphism{x.entry}
}

type BuilderMapMorphism struct {
	entry *AtlasEntry
}

func (x *BuilderMapMorphism) Complete() *AtlasEntry {
	return x.entry
}

func (x *BuilderMapMorphism) SetKeySortMode(km KeySortMode) *BuilderMapMorphism {
	switch km {
	case KeySortMode_Default, KeySortMode_RFC7049:
		x.entry.MapMorphism.KeySortMode = km
	default:
		panic(fmt.Errorf("invalid key sort mode %q", km))
	}
	return x
}
