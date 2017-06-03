package refmt

import (
	"bytes"
	"io"

	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/obj"
	"github.com/polydawn/refmt/obj/atlas"
)

func NewCborEncoder(wr io.Writer) *CborEncoder {
	return NewAtlasedCborEncoder(wr, atlas.MustBuild())
}

func NewAtlasedCborEncoder(wr io.Writer, atl atlas.Atlas) *CborEncoder {
	enc := &CborEncoder{
		marshaller: obj.NewMarshaler(atl),
		encoder:    cbor.NewEncoder(wr),
	}
	enc.pump = TokenPump{
		enc.marshaller,
		enc.encoder,
	}
	return enc
}

type CborEncoder struct {
	marshaller *obj.MarshalDriver
	encoder    *cbor.Encoder
	pump       TokenPump
}

func (d *CborEncoder) Marshal(v interface{}) error {
	d.marshaller.Bind(v)
	d.encoder.Reset()
	return d.pump.Run()
}

func CborEncode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewCborEncoder(&buf).Marshal(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func CborEncodeAtlased(atl atlas.Atlas, v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewAtlasedCborEncoder(&buf, atl).Marshal(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewCborDecoder(r io.Reader) *CborDecoder {
	return &CborDecoder{r}
}

type CborDecoder struct {
	r io.Reader
}

func (d *CborDecoder) Unmarshal(v interface{}) error {
	return TokenPump{
		cbor.NewDecoder(d.r),
		obj.NewUnmarshaler(atlas.MustBuild()),
	}.Run()
}

func CborDecode(v interface{}, b []byte) error {
	return NewCborDecoder(bytes.NewBuffer(b)).Unmarshal(v)
}
