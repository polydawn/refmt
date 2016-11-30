package xlate

import (
	"io"

	"github.com/polydawn/go-xlate/obj"
)

func NewJsonEncoder(wr io.Writer) *JsonEncoder {
	return &JsonEncoder{wr}
}

type JsonEncoder struct {
	wr io.Writer
}

func (d *JsonEncoder) Marshal(v interface{}) {
	TokenPump{
		nil, // todo get the rest of obj.NewMarshaller(v) in place
		nil, // todo get the whole json package in place
	}.Run()
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
