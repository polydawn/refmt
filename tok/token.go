package tok

/*
	Token is either one of magic const tokens (used to denote beginning and
	ending of maps and arrays), or an address to a primitive (string, int, etc).
*/
type Token interface{}

var (
	Token_MapOpen  Token = '{'
	Token_MapClose Token = '}'
	Token_ArrOpen  Token = '['
	Token_ArrClose Token = ']'
)
