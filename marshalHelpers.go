package refmt

import (
	"io"

	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj/atlas"
)

type EncodeOptions interface {
	IsEncodeOptions() // marker method.
}

func Marshal(opts EncodeOptions, v interface{}) ([]byte, error) {
	switch opts.(type) {
	case json.EncodeOptions:
		return json.Marshal(v)
	case cbor.EncodeOptions:
		return cbor.Marshal(v)
	default:
		panic("incorrect usage: unknown EncodeOptions type")
	}
}

func MarshalAtlased(opts EncodeOptions, v interface{}, atl atlas.Atlas) ([]byte, error) {
	switch opts.(type) {
	case json.EncodeOptions:
		return json.MarshalAtlased(v, atl)
	case cbor.EncodeOptions:
		return cbor.MarshalAtlased(v, atl)
	default:
		panic("incorrect usage: unknown EncodeOptions type")
	}
}

type Marshaller interface {
	Marshal(v interface{}) error
}

func NewMarshaller(opts EncodeOptions, wr io.Writer) Marshaller {
	switch opts.(type) {
	case json.EncodeOptions:
		return json.NewMarshaller(wr)
	case cbor.EncodeOptions:
		return cbor.NewMarshaller(wr)
	default:
		panic("incorrect usage: unknown EncodeOptions type")
	}
}

func NewMarshallerAtlased(opts EncodeOptions, wr io.Writer, atl atlas.Atlas) Marshaller {
	switch opts.(type) {
	case json.EncodeOptions:
		return json.NewMarshallerAtlased(wr, atl)
	case cbor.EncodeOptions:
		return cbor.NewMarshallerAtlased(wr, atl)
	default:
		panic("incorrect usage: unknown EncodeOptions type")
	}
}
