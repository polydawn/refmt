package xlate

/*
	Returns a new `Destination` implementation that will write the value(s)
	it receives into a the given `interface{}` reference.
*/
func NewVarDestination(where interface{}) Destination {
	return &varDest{where}
}

type varDest struct {
	Var interface{}
}

func (d *varDest) OpenMap()                              {}
func (d *varDest) WriteMapEntry(k string, v interface{}) {}
func (d *varDest) CloseMap()                             {}

func (d *varDest) OpenArray()                    {}
func (d *varDest) WriteArrayEntry(v interface{}) {}
func (d *varDest) CloseArray()                   {}

func (d *varDest) WriteString(v string) {
	(*(d.Var.(*string))) = v
}
