package cbor

type EncodeOptions struct {
	// there aren't a ton of options for cbor, but we still need this
	// for use as a sigil for the top-level refmt methods to demux on.
}

// marker method -- you may use this type to instruct `refmt.Marshal`
// what kind of encoder to use.
func (EncodeOptions) IsEncodeOptions() {}

type DecodeOptions struct {
	CoerceUndefToNull bool

	// RejectIndefinite causes the decoder to error when it encounters an
	// indefinite-length encoding (major types 2, 3, 4 or 5 with the
	// 0x1f indefinite indicator). The error is returned as soon as the
	// indefinite-length sigil byte is seen, before any chunks are read or
	// allocated. Useful for codecs that forbid indefinite-length values
	// (notably DAG-CBOR).
	RejectIndefinite bool

	// MaxIndefiniteSize caps the cumulative size, in bytes, of an
	// indefinite-length bytes or string value during chunk aggregation.
	// Decoding errors when the running total would exceed this. When zero,
	// a default of 32 MiB is used (matching the per-chunk maximum, so
	// indefinite values can't grow larger in total than a single
	// definite-length value).
	MaxIndefiniteSize int

	// RejectNonMinimalInteger rejects CBOR heads whose integer argument is
	// encoded in more bytes than necessary. Applies to uints, negative
	// ints, length headers (bytes/strings/arrays/maps) and tag headers.
	// Required by codecs that mandate minimal encoding (e.g. DAG-CBOR).
	RejectNonMinimalInteger bool

	// RejectNaN causes the decoder to error when a float value decodes to
	// NaN (any of its many bit representations). Required by codecs that
	// forbid NaN, including DAG-CBOR.
	RejectNaN bool

	// RejectInfinity causes the decoder to error when a float value
	// decodes to +Inf or -Inf. Required by codecs that forbid infinities,
	// including DAG-CBOR.
	RejectInfinity bool

	// RejectNarrowFloat rejects float values encoded as 16-bit (0xf9) or 32-bit
	// (0xfa) at the sigil byte, before any payload is read. Required by
	// codecs that mandate float64-only encoding (DAG-CBOR).
	RejectNarrowFloat bool

	// future: options to validate canonical serial order
}

// defaultMaxIndefiniteSize matches the existing per-chunk size cap in
// decodeBytesOrStringIndefinite (33554432 bytes / 32 MiB). It bounds the
// total accumulator so an indefinite-length value cannot grow larger than
// what a single definite-length value would have been allowed to.
const defaultMaxIndefiniteSize = 33554432

func (cfg DecodeOptions) maxIndefiniteSize() int {
	if cfg.MaxIndefiniteSize > 0 {
		return cfg.MaxIndefiniteSize
	}
	return defaultMaxIndefiniteSize
}

// marker method -- you may use this type to instruct `refmt.Marshal`
// what kind of encoder to use.
func (DecodeOptions) IsDecodeOptions() {}
