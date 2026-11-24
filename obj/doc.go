/*
The `obj` package defines Marshaller and Unmarshaller types,
which can be used to convert in-memory values to token streams,
and token streams to unpack in-memory values.

When decoding an integer token into a concrete Go integer type, Unmarshaller
returns ErrUnmarshalTypeCantFit if the value is outside that type's range.
*/
package obj
