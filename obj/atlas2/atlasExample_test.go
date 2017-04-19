// +build none
// this is only a syntax playground file at the moment

package atlas_test

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

func example() {
	x := atlas.Build(
		atlas.BuildEntry(typeHintObj{}).StructMap().
			AddField("FieldName", StructMapEntry{SerialName: "fn", OmitEmpty: true}).
			AddField("Nested.Thing", StructMapEntry{SerialName: "nt"}).
			Complete(),
		// and carry on; this `Build` method takes `AtlasEntry...`.
	)
}
