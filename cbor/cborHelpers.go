package cbor

import (
	"bytes"
	"io"

	"github.com/polydawn/refmt/obj"
	"github.com/polydawn/refmt/obj/atlas"
	"github.com/polydawn/refmt/shared"
)

// All of the methods in this file are exported,
// and their names and type declarations are intended to be
// identical to the naming and types of the golang stdlib
// 'encoding/json' packages, with ONE EXCEPTION:
// what stdlib calls "NewEncoder", we call "NewMarshaler";
// what stdlib calls "NewDecoder", we call "NewUnmarshaler".
// You should be able to migrate with a sed script!
//
// (In refmt, the encoder/decoder systems are for token streams;
// if you're talking about object mapping, we consistently
// refer to that as marshaling/unmarshaling.)
//
// Most methods also have an "Atlased" variant,
// which lets you specify advanced type mapping instructions.

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewMarshaler(&buf).Marshal(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Marshaler struct {
	marshaler *obj.Marshaler
	encoder   *Encoder
	pump      shared.TokenPump
}

func (x *Marshaler) Marshal(v interface{}) error {
	x.marshaler.Bind(v)
	x.encoder.Reset()
	return x.pump.Run()
}

func NewMarshaler(wr io.Writer) *Marshaler {
	return NewMarshalerAtlased(wr, atlas.MustBuild())
}

func NewMarshalerAtlased(wr io.Writer, atl atlas.Atlas) *Marshaler {
	x := &Marshaler{
		marshaler: obj.NewMarshaler(atl),
		encoder:   NewEncoder(wr),
	}
	x.pump = shared.TokenPump{
		x.marshaler,
		x.encoder,
	}
	return x
}

func Unmarshal(data []byte, v interface{}) error {
	return NewUnmarshaler(bytes.NewBuffer(data)).Unmarshal(v)
}

type Unmarshaler struct {
	unmarshaler *obj.Unmarshaler
	decoder     *Decoder
	pump        shared.TokenPump
}

func (x *Unmarshaler) Unmarshal(v interface{}) error {
	x.unmarshaler.Bind(v)
	x.decoder.Reset()
	return x.pump.Run()
}

func NewUnmarshaler(r io.Reader) *Unmarshaler {
	return NewUnmarshalerAtlased(r, atlas.MustBuild())
}
func NewUnmarshalerAtlased(r io.Reader, atl atlas.Atlas) *Unmarshaler {
	x := &Unmarshaler{
		unmarshaler: obj.NewUnmarshaler(atl),
		decoder:     NewDecoder(r),
	}
	x.pump = shared.TokenPump{
		x.unmarshaler,
		x.decoder,
	}
	return x
}
