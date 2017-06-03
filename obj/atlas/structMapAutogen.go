package atlas

import (
	"reflect"
	"strings"
	"unicode"
)

func AutogenerateStructMapEntry(rt reflect.Type) *AtlasEntry {
	return AutogenerateStructMapEntryUsingTags(rt, "refmt")
}

func AutogenerateStructMapEntryUsingTags(rt reflect.Type, tagName string) *AtlasEntry {
	entry := &AtlasEntry{
		Type:      rt,
		StructMap: &StructMap{Fields: exploreFields(rt, tagName)},
	}
	return entry
}

// exploreFields returns a list of fields that StructAtlas should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func exploreFields(rt reflect.Type, tagName string) []StructMapEntry {
	// Anonymous fields to explore at the current level and the next.
	current := []StructMapEntry{}
	next := []StructMapEntry{{Type: rt}}

	// Count of queued names for current level and the next.
	count := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []StructMapEntry

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.Type] {
				continue
			}
			visited[f.Type] = true

			// Scan f.Type for fields to include.
			for i := 0; i < f.Type.NumField(); i++ {
				sf := f.Type.Field(i)
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
				route := make([]int, len(f.ReflectRoute)+1)
				copy(route, f.ReflectRoute)
				route[len(f.ReflectRoute)] = i

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
					fields = append(fields, StructMapEntry{
						SerialName:   name, // TODO default to downcaseing
						ReflectRoute: route,
						Type:         ft,
						tagged:       tagged,
						OmitEmpty:    opts.Contains("omitempty"),
					})
					if count[f.Type] > 1 {
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
					next = append(next, StructMapEntry{
						ReflectRoute: route,
						Type:         ft,
					})
				}
			}
		}
	}

	// TODO get the sorting and annihilation in here

	return fields
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
