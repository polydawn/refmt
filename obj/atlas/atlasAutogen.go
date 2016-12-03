package atlas

import (
	"reflect"
	"sort"
	"strings"
	"unicode"
)

func GenerateAtlas(t reflect.Type) *Atlas {
	return GenerateAtlasUsingTags(t, "atlas")
}

func GenerateAtlasUsingTags(t reflect.Type, tagName string) *Atlas {
	fields := exploreFields(t, tagName)
	atl := &Atlas{
		Fields: make([]Entry, len(fields)),
	}
	for i, v := range fields {
		atl.Fields[i] = v.Entry
	}
	return atl
}

// A field represents a single field found in a struct.
type atlasGenField struct {
	Entry // we're populating this

	// the rest will be forgotten after exploration:

	tag bool         // if named by tag -- used for export prio
	typ reflect.Type // convenient handle
}

// atlasGenField_byName sorts field by name,
// breaking ties with depth,
// then breaking ties with "name came from tag",
// then breaking ties with FieldRoute sequence.
type atlasGenField_byName []atlasGenField

func (x atlasGenField_byName) Len() int { return len(x) }

func (x atlasGenField_byName) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x atlasGenField_byName) Less(i, j int) bool {
	if x[i].Name != x[j].Name {
		return x[i].Name < x[j].Name
	}
	if len(x[i].FieldRoute) != len(x[j].FieldRoute) {
		return len(x[i].FieldRoute) < len(x[j].FieldRoute)
	}
	if x[i].tag != x[j].tag {
		return x[i].tag
	}
	return atlasGenField_byFieldRoute(x).Less(i, j)
}

// atlasGenField_byFieldRoute sorts field by FieldRoute sequence
// (e.g., roughly source declaration order within each type).
type atlasGenField_byFieldRoute []atlasGenField

func (x atlasGenField_byFieldRoute) Len() int { return len(x) }

func (x atlasGenField_byFieldRoute) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x atlasGenField_byFieldRoute) Less(i, j int) bool {
	for k, xik := range x[i].FieldRoute {
		if k >= len(x[j].FieldRoute) {
			return false
		}
		if xik != x[j].FieldRoute[k] {
			return xik < x[j].FieldRoute[k]
		}
	}
	return len(x[i].FieldRoute) < len(x[j].FieldRoute)
}

// exploreFields returns a list of fields that StructAtlas should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func exploreFields(t reflect.Type, tagName string) []atlasGenField {
	// Anonymous fields to explore at the current level and the next.
	current := []atlasGenField{}
	next := []atlasGenField{{typ: t}}

	// Count of queued names for current level and the next.
	count := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []atlasGenField

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.PkgPath != "" && !sf.Anonymous { // unexported
					continue
				}
				tag := sf.Tag.Get(tagName)
				if tag == "-" {
					continue
				}
				name, opts := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}
				route := make([]int, len(f.FieldRoute)+1)
				copy(route, f.FieldRoute)
				route[len(f.FieldRoute)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					fields = append(fields, atlasGenField{
						Entry: Entry{
							Name:       name,
							FieldRoute: route,
							OmitEmpty:  opts.Contains("omitempty"),
						},
						tag: tagged,
						typ: ft,
					})
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, atlasGenField{
						Entry: Entry{
							Name:       ft.Name(),
							FieldRoute: route,
						},
						typ: ft,
					})
				}
			}
		}
	}

	sort.Sort(atlasGenField_byName(fields))

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with JSON tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field index length. Loop over names; for each name, delete
	// hidden fields by choosing the one dominant field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.Name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.Name != name {
				break
			}
		}
		if advance == 1 { // Only one field with this name
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(atlasGenField_byFieldRoute(fields))

	return fields
}

// dominantField looks through the fields, all of which are known to
// have the same name, to find the single field that dominates the
// others using Go's embedding rules, modified by the presence of
// JSON tags. If there are multiple top-level fields, the boolean
// will be false: This condition is an error in Go and we skip all
// the fields.
func dominantField(fields []atlasGenField) (atlasGenField, bool) {
	// The fields are sorted in increasing index-length order. The winner
	// must therefore be one with the shortest index length. Drop all
	// longer entries, which is easy: just truncate the slice.
	length := len(fields[0].FieldRoute)
	tagged := -1 // Index of first tagged field.
	for i, f := range fields {
		if len(f.FieldRoute) > length {
			fields = fields[:i]
			break
		}
		if f.tag {
			if tagged >= 0 {
				// Multiple tagged fields at the same level: conflict.
				// Return no field.
				return atlasGenField{}, false
			}
			tagged = i
		}
	}
	if tagged >= 0 {
		return fields[tagged], true
	}
	// All remaining fields have the same length. If there's more than one,
	// we have a conflict (two fields named "X" at the same level) and we
	// return no field.
	if len(fields) > 1 {
		return atlasGenField{}, false
	}
	return fields[0], true
}

func fieldByIndex(v reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}

// tagOptions is the string following a comma in a struct field's
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}
