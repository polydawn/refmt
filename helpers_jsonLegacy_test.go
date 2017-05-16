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
	enc := &JsonLegacyEncoder{
		marshaller: objLegacy.NewMarshaler(&objLegacy.Suite{}),
		serializer: json.NewSerializer(wr),
	}
	enc.pump = TokenPump{
		enc.marshaller,
		enc.serializer,
	}
	return enc
}

type JsonLegacyEncoder struct {
	marshaller *objLegacy.MarshalDriver
	serializer *json.Serializer
	pump       TokenPump
}

func (d *JsonLegacyEncoder) Marshal(v interface{}) error {
	d.marshaller.Bind(v)
	d.serializer.Reset()
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
