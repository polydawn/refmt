package refmt

import (
	"bytes"
	"io"

	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj"
	"github.com/polydawn/refmt/obj/atlas"
)

func NewJsonEncoder(wr io.Writer) *JsonEncoder {
	return NewAtlasedJsonEncoder(wr, atlas.MustBuild())
}

func NewAtlasedJsonEncoder(wr io.Writer, atl atlas.Atlas) *JsonEncoder {
	enc := &JsonEncoder{
		marshaller: obj.NewMarshaler(atl),
		encoder:    json.NewEncoder(wr),
	}
	enc.pump = TokenPump{
		enc.marshaller,
		enc.encoder,
	}
	return enc
}

type JsonEncoder struct {
	marshaller *obj.MarshalDriver
	encoder    *json.Encoder
	pump       TokenPump
}

func (d *JsonEncoder) Marshal(v interface{}) error {
	d.marshaller.Bind(v)
	d.encoder.Reset()
	return d.pump.Run()
}

func JsonEncode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewJsonEncoder(&buf).Marshal(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func JsonEncodeAtlased(v interface{}, atl atlas.Atlas) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewAtlasedJsonEncoder(&buf, atl).Marshal(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewJsonDecoder(r io.Reader) *JsonDecoder {
	return &JsonDecoder{r}
}

type JsonDecoder struct {
	r io.Reader
}

func (d *JsonDecoder) Unmarshal(v interface{}) {
	TokenPump{
		nil, // todo get the whole json package in place
		obj.NewUnmarshaler(atlas.MustBuild()),
	}.Run()
}
