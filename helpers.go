package xlate

import (
	"bytes"
	"io"

	"github.com/polydawn/go-xlate/json"
	"github.com/polydawn/go-xlate/obj"
)

func NewJsonEncoder(wr io.Writer) *JsonEncoder {
	return &JsonEncoder{wr}
}

type JsonEncoder struct {
	wr io.Writer
}

func (d *JsonEncoder) Marshal(v interface{}) error {
	return TokenPump{
		obj.NewMarshaler(&obj.Suite{}, v),
		json.NewSerializer(d.wr),
	}.Run()
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
		obj.NewUnmarshaler(v),
	}.Run()
}
