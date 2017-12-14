package atlas_test

import (
	"."
)

// Design notes:
//
// - Remember that we do want the AtlasEntry to be definable with each type.
//   Locally.  There's no reason you shouldn't be able to define them
//   physically adjacent to the source of the structures in your source code.
//
// - We also don't want to repeat ourselves (or more precisely, we absolutely
//   shouldn't make the user repeat themselves, but we also don't want to
//   repeat ourselves in docs and factory methods unnecessarily), so,
//   given that keeping a reflect.Type instance on hand appears to be common,
//   let's start building any builder with that.  Then, diversify.

func ExampleAtlasBuilding() {
	type typeExample1 struct {
		FieldName string
		Nested    struct {
			Thing int
		}
	}

	atl, err := atlas.Build(
		atlas.BuildEntry(typeExample1{}).StructMap().
			AddField("FieldName", atlas.StructMapEntry{SerialName: "fn", OmitEmpty: true}).
			AddField("Nested.Thing", atlas.StructMapEntry{SerialName: "nt"}).
			Complete(),
		atlas.BuildEntry(map[string]typeExample1{}).MapMorphism().
			SetKeySortMode(atlas.KeySortMode_RFC7049).
			Complete(),
		// and carry on; this `Build` method takes `AtlasEntry...` as a vararg.
	)
	_, _ = atl, err

	// Output:
}
