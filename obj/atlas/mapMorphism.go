package atlas

// A type to enumerate key sorting modes.
type KeySortMode string

const (
	KeySortMode_Default = KeySortMode("default") // e.g. lexical string sort for strings, etc
	KeySortMode_RFC7049 = KeySortMode("rfc7049") // "Canonical" as proposed by rfc7049 § 3.9 (shorter byte sequences sort to top)
)

type MapMorphism struct {
	KeySortMode KeySortMode
}
