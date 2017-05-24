// This file is in the '_test' namespace because we want to use these helpers in benchmarks,
// but do not wish to expose them to users or publish them in the API docs anymore.

package refmt

import (
	"bytes"
	"io"

	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/objLegacy"
)

func NewJsonLegacyEncoder(wr io.Writer) *JsonLegacyEncoder {
	return NewAtlasedJsonLegacyEncoder(wr, &objLegacy.Suite{})
}

func NewAtlasedJsonLegacyEncoder(wr io.Writer, suite *objLegacy.Suite) *JsonLegacyEncoder {
	enc := &JsonLegacyEncoder{
		marshaller: objLegacy.NewMarshaler(suite),
		encoder:    json.NewEncoder(wr),
	}
	enc.pump = TokenPump{
		enc.marshaller,
		enc.encoder,
	}
	return enc
}

type JsonLegacyEncoder struct {
	marshaller *objLegacy.MarshalDriver
	encoder    *json.Encoder
	pump       TokenPump
}

func (d *JsonLegacyEncoder) Marshal(v interface{}) error {
	d.marshaller.Bind(v)
	d.encoder.Reset()
	return d.pump.Run()
}

func JsonLegacyEncode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewJsonLegacyEncoder(&buf).Marshal(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewJsonLegacyDecoder(r io.Reader) *JsonLegacyDecoder {
	return &JsonLegacyDecoder{r}
}

type JsonLegacyDecoder struct {
	r io.Reader
}

func (d *JsonLegacyDecoder) Unmarshal(v interface{}) {
	TokenPump{
		nil, // todo get the whole json package in place
		objLegacy.NewUnmarshaler(v),
	}.Run()
}
