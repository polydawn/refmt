package atlas

type ErrEntryInvalid struct {
	Msg string
}

func (e ErrEntryInvalid) Error() string {
	return "atlas.Entry invalid: " + e.Msg
}
