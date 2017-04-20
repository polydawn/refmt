package obj

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflectInvariants(t *testing.T) {
	describe := func(thing interface{}) {
		rt := reflect.TypeOf(thing)
		fmt.Printf("type: %#v (%q)\n", rt, rt)
		fmt.Printf("kind: %#v (%q)\n", rt.Kind(), rt.Kind())
		fmt.Printf("name: %q\n", rt.Name())
		fmt.Printf("PkgPath: %q\n", rt.PkgPath())
	}
	describe("")
	fmt.Println()
	type StrTypedef string
	describe(StrTypedef(""))
	fmt.Println()

	// I guess the cheapest and sanest way to check if something is one of the
	// builtin primitives is just checking the rtid pointer equality outright?
	//
	// I think `PkgPath() == "" && (... kind is not slice, etc...)` might also work,
	// but the documentation doesn't make it very clear if that's an intended use
	// of the PkgPath function, and this seems much harder and less reliable
	// than the "hack".

	fmt.Printf("rtid: %v\n", reflect.ValueOf(reflect.TypeOf("")).Pointer())
	fmt.Printf("rtid: %v\n", reflect.ValueOf(reflect.TypeOf("")).Pointer())
	fmt.Printf("rtid: %v\n", reflect.ValueOf(reflect.TypeOf(StrTypedef(""))).Pointer())
	fmt.Printf("rtid: %v\n", reflect.ValueOf(reflect.TypeOf(StrTypedef(""))).Pointer())
}
