package atlas

// Error type raised when an atlas.Entry is invalid, missing required values,
// or otherwise extremely wrong.
type ErrEntryInvalid struct {
	Msg string
}

func (e ErrEntryInvalid) Error() string {
	return "atlas.Entry invalid: " + e.Msg
}

// Error type raised when initializing an Atlas, and field entries do
// not resolve against the type.
// (If you recently refactored names of fields in your types, check
// to make sure you updated any references to those fields by name to match!)
type ErrStructureMismatch struct {
	TypeName string
	Reason   string
}

func (e ErrStructureMismatch) Error() string {
	return "atlas does not match type: " + e.TypeName + " " + e.Reason
}
