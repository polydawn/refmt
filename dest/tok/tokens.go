package tok

import (
	"bytes"
	"fmt"
)

type TokenKind int

const (
	TokenKind_Invalid TokenKind = iota
	TokenKind_OpenMap
	TokenKind_MapKey
	TokenKind_CloseMap
	TokenKind_OpenArray
	TokenKind_CloseArray
	TokenKind_ValNull
	TokenKind_ValString
)

func (tk TokenKind) String() string {
	switch tk {
	case TokenKind_Invalid:
		return "INVALID"
	case TokenKind_OpenMap:
		return "{"
	case TokenKind_MapKey:
		return "k"
	case TokenKind_CloseMap:
		return "}"
	case TokenKind_OpenArray:
		return "["
	case TokenKind_CloseArray:
		return "]"
	case TokenKind_ValNull:
		return "-"
	case TokenKind_ValString:
		return "s"
	}
	panic("missing case")
}

type Token struct {
	Kind TokenKind
	Data interface{}
}

func (t Token) String() string {
	if t.Data == nil {
		return fmt.Sprintf("<%s>", t.Kind)
	}
	return fmt.Sprintf("<%s:%s>", t.Kind, t.Data)
}

type Tokens []Token

func (ts Tokens) String() string {
	var buf bytes.Buffer
	buf.WriteString("Tokens[\n")
	for _, t := range ts {
		buf.WriteByte('\t')
		fmt.Fprintf(&buf, "%s", t)
		buf.WriteByte('\n')
	}
	buf.WriteString("]")
	return buf.String()
}

func (ts Tokens) GoString() string {
	return ts.String()
}
