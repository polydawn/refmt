package refmt

import (
	"bytes"
	"io"

	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj"
	"github.com/polydawn/refmt/obj/atlas"
)

func NewJsonEncoder(wr io.Writer) *JsonEncoder {
	enc := &JsonEncoder{
		marshaller: obj.NewMarshaler(atlas.MustBuild()),
		serializer: json.NewSerializer(wr),
	}
	enc.pump = TokenPump{
		enc.marshaller,
		enc.serializer,
	}
	return enc
}

type JsonEncoder struct {
	marshaller *obj.MarshalDriver
	serializer *json.Serializer
	pump       TokenPump
}

func (d *JsonEncoder) Marshal(v interface{}) error {
	d.marshaller.Bind(v)
	d.serializer.Reset()
	return d.pump.Run()
}

func JsonEncode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewJsonEncoder(&buf).Marshal(v); err != nil {
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
