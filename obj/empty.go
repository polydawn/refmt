package obj

import "reflect"

// The missing definition of 'reflect.IsZero' you've always wanted.
//
// This definition always considers structs non-empty
// (checking struct emptiness can be a relatively costly operation, O(size_of_struct);
// all other kinds can define emptiness in constant time).
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// zeroishness checks on a struct are possible in theory but we're passing for now.
// There's no easy route exposed from the reflect package, due to a variety of reasons:
//
//   - The `DeepEqual` method would almost suffice, but we'd need one that takes
//     `reflect.Value` instead of `interface{}` params, because we already have
//     the former, and un/re-boxing those into two `interface{}` values is two
//     heap mallocs and that's a wildly unacceptable performance overhead for this.
//     `DeepEqual` calls `deepValueEqual`, which does what we want, but...
//   - We can't easily copy the `deepValueEqual` method.  It imports the unsafe
//     package.  This is undesirable because it would reduce the portability of refmt.
//
// It's possible we could produce a struct isZero method which is *simpler* than
// `deepValueEqual`, because we can actually just halt on any non-zero pointer,
// and without any need for following pointers we also need no cycle detection.
// But this is an exercise I'm leaving for later.  PRs welcome if someone wants it.
