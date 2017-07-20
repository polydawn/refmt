// This file is in the '_test' namespace because we want to use these helpers in benchmarks,
// but do not wish to expose them to users or publish them in the API docs anymore.

package refmt

import (
	"bytes"
	"io"

	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/objLegacy"
)

func NewCborLegacyEncoder(wr io.Writer) *CborLegacyEncoder {
	return NewAtlasedCborLegacyEncoder(wr, &objLegacy.Suite{})
}

func NewAtlasedCborLegacyEncoder(wr io.Writer, suite *objLegacy.Suite) *CborLegacyEncoder {
	enc := &CborLegacyEncoder{
		marshaller: objLegacy.NewMarshaller(suite),
		encoder:    cbor.NewEncoder(wr),
	}
	enc.pump = TokenPump{
		enc.marshaller,
		enc.encoder,
	}
	return enc
}

type CborLegacyEncoder struct {
	marshaller *objLegacy.MarshalDriver
	encoder    *cbor.Encoder
	pump       TokenPump
}

func (d *CborLegacyEncoder) Marshal(v interface{}) error {
	d.marshaller.Bind(v)
	d.encoder.Reset()
	return d.pump.Run()
}

func CborLegacyEncode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewCborLegacyEncoder(&buf).Marshal(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewCborLegacyDecoder(r io.Reader) *CborLegacyDecoder {
	return &CborLegacyDecoder{r}
}

type CborLegacyDecoder struct {
	r io.Reader
}

func (d *CborLegacyDecoder) Unmarshal(v interface{}) error {
	return TokenPump{
		cbor.NewDecoder(d.r),
		objLegacy.NewUnmarshaller(&objLegacy.Suite{}),
	}.Run()
}

func CborLegacyDecode(v interface{}, b []byte) error {
	return NewCborLegacyDecoder(bytes.NewBuffer(b)).Unmarshal(v)
}
