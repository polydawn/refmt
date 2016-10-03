/*
	Implements a really trivial destination that buffers the entire stream of
	writes it ever receives -- useful for debugging and testing.
*/
package tok

import (
	"github.com/polydawn/go-xlate"
)

var _ xlate.Destination = &Destination{}

/*
	Returns a new `Destination` implementation that will buffer a record of
	all the calls it receives.
*/
func NewDestination() *Destination {
	return &Destination{}
}

type Destination struct {
	Tokens Tokens
}

func (d *Destination) push(tk TokenKind, data interface{}) {
	d.Tokens = append(d.Tokens, Token{tk, data})
}

func (d *Destination) OpenMap()             { d.push(TokenKind_OpenMap, nil) }
func (d *Destination) WriteMapKey(k string) { d.push(TokenKind_MapKey, k) }
func (d *Destination) CloseMap()            { d.push(TokenKind_CloseMap, nil) }
func (d *Destination) OpenArray()           { d.push(TokenKind_OpenArray, nil) }
func (d *Destination) CloseArray()          { d.push(TokenKind_CloseArray, nil) }
func (d *Destination) WriteNull()           { d.push(TokenKind_ValNull, nil) }
func (d *Destination) WriteString(v string) { d.push(TokenKind_ValString, v) }
