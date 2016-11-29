go-xlate
========

[![GoDoc](https://godoc.org/github.com/polydawn/go-xlate?status.svg)](https://godoc.org/github.com/polydawn/go-xlate)



Why?
----

Mostly because I have some types which I need to encode in two different ways, and that needs to not suck, and that totally sucks with most serialization libraries I've used.

More broadly, I want a single library that can handle my serialization -- with the possibility of different setups on the same types -- and if it can do general object traversals, e.g. a deepcopy, that also seems like... just something that should be natural.

So it seems like there should be some way to define object walkers... and some way to define emitting a stream of values during such a walk... and they should be pretty separate.

Thusly was this library thrust into the world: `xlate.Mapper` to define the walks, and `xlate.Destination` to marshal the value stream.



Show me.
--------

```go
// This is an example struct.
type AA struct{}

// This is a MapperFunc
func Map_Wildcard_toString(dest Destination, input interface{}) {
	dest.WriteString(fmt.Sprintf("%s", input))
}

// This constructs a mapper.
mapper := NewMapper(MapperSetup{
	{AA{}, Map_Wildcard_toStringOfType},
})

// This constructs a destination (to just put things back in memory as an object).
var result string
destination := NewVarDestination(&result)

// This has the mapper walk over an object and stream it to the destination!
mapper.Map(AA{}, destination)

// `result` now contains "AA"!
```

The key bit here is:

- you can swap out those functions in the mapper setup
- and that works for serialization modes even deep in other structures;
- at the same time, you can swap out that destination for a json encoder
- and the mapper part *doesn't care* -- no complex interactions between the layers.
