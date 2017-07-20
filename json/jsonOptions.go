package json

type EncodeOptions struct {
	// future: indentation opts and such
}

// marker method -- you may use this type to instruct `refmt.Marshal`
// what kind of encoder to use.
func (EncodeOptions) IsEncodeOptions() {}

type DecodeOptions struct {
	// future: options to validate canonical serial order
}

// marker method -- you may use this type to instruct `refmt.Marshal`
// what kind of encoder to use.
func (DecodeOptions) IsDecodeOptions() {}
