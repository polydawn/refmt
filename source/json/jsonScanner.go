package json

import "strconv"

/*
	notes:

	- `scanTransition_BeginLiteral` is used WAY too heavily: trues, falses, ints, floats, strings, all there.
	  - i want to emit those tokens all separately (i think!) and it's a little funny to me that the scanner wouldn't be the one to be clear about that.

*/

// A SyntaxError is a description of a JSON syntax error.
type SyntaxError struct {
	msg    string // description of error
	Offset int64  // error occurred after reading Offset bytes
	// TODO i'd like to keep lines and offset within a line here as well
}

func (e *SyntaxError) Error() string { return e.msg }

type scanTransition int8

// These values are returned by the state transition functions
// assigned to scanner.state and the method scanner.eof.
// They give details about the current state of the scan that
// callers might be interested to know about.
// It is okay to ignore the return value of any particular
// call to scanner.state: if one call returns scanTransition_Error,
// every subsequent call will return scanTransition_Error too.
const (
	// Continue.
	scanTransition_Continue     = iota // uninteresting byte
	scanTransition_BeginLiteral        // end implied by next result != scanTransition_Continue
	scanTransition_BeginMap            // begin map
	scanTransition_MapKey              // just finished map key (string)
	scanTransition_MapValue            // just finished non-last map value
	scanTransition_EndMap              // end map (implies scanTransition_MapValue if possible)
	scanTransition_BeginArray          // begin array
	scanTransition_ArrayValue          // just finished array value
	scanTransition_EndArray            // end array (implies scanTransition_ArrayValue if possible)
	scanTransition_SkipSpace           // space byte; can skip; known to be last "continue" result

	// Stop.
	scanTransition_End   // top-level value ended *before* this byte; known to be first "stop" result
	scanTransition_Error // hit an error, scanner.err.
)

type parsePhase int8

// These values are stored in the parseState stack.
// They give the current state of a composite value
// being scanned. If the parser is inside a nested value
// the parseState describes the nested state, outermost at entry 0.
const (
	parsePhase_MapKey     = iota // parsing map key (before colon)
	parsePhase_MapValue          // parsing map value (after colon)
	parsePhase_ArrayValue        // parsing array value
)

/*
	The scanner state machine advances one byte at a time,
	yielding a signal code each time it finds an edge in the content.

	The scanner does not ensure that the tokens are in a semantically valid order, etc.
*/
type scanner struct {
	// The step is a func to be called to execute the next transition.
	step func(*scanner, byte) scanTransition

	// Reached end of top-level value.
	endTop bool // TODO REVIEW this is only ever *set*? // think it's supposed be used by an eof state that we haven't seemed to need yet.

	// Stack of what we're in the middle of - array values, map keys, map values.
	parseState []int

	// Error that happened, if any.
	err error

	// 1-byte redo (see undo method)
	redo      bool
	redoByte  byte
	redoCode  scanTransition
	redoState func(*scanner, byte) scanTransition

	// total bytes consumed, updated by decoder.Decode
	bytes int64 // FIXME again, srsly, why is this cross-reaching?
}

func (s *scanner) reset() {
	s.step = stateBeginValue
	s.parseState = s.parseState[0:0]
}

// pushParseState pushes a new parse state p onto the parse stack.
func (s *scanner) pushParseState(p int) {
	s.parseState = append(s.parseState, p)
}

// popParseState pops a parse state (already obtained) off the stack
// and updates s.step accordingly.
func (s *scanner) popParseState() {
	n := len(s.parseState) - 1
	s.parseState = s.parseState[0:n]
	s.redo = false
	if n == 0 {
		s.step = stateEndTop
		s.endTop = true
	} else {
		s.step = stateEndValue
	}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

// stateBeginValueOrEmpty is the state after reading `[`.
func stateBeginValueOrEmpty(s *scanner, c byte) scanTransition {
	if c <= ' ' && isSpace(c) {
		return scanTransition_SkipSpace
	}
	if c == ']' {
		return stateEndValue(s, c)
	}
	return stateBeginValue(s, c)
}

// stateBeginValue is the state at the beginning of the input.
func stateBeginValue(s *scanner, c byte) scanTransition {
	if c <= ' ' && isSpace(c) {
		return scanTransition_SkipSpace
	}
	switch c {
	case '{':
		s.step = stateBeginStringOrEmpty
		s.pushParseState(parsePhase_MapKey)
		return scanTransition_BeginMap
	case '[':
		s.step = stateBeginValueOrEmpty
		s.pushParseState(parsePhase_ArrayValue)
		return scanTransition_BeginArray
	case '"':
		s.step = stateInString
		return scanTransition_BeginLiteral
	case '-':
		s.step = stateNeg
		return scanTransition_BeginLiteral
	case '0': // beginning of 0.123
		s.step = state0
		return scanTransition_BeginLiteral
	case 't': // beginning of true
		s.step = stateT
		return scanTransition_BeginLiteral
	case 'f': // beginning of false
		s.step = stateF
		return scanTransition_BeginLiteral
	case 'n': // beginning of null
		s.step = stateN
		return scanTransition_BeginLiteral
	}
	if '1' <= c && c <= '9' { // beginning of 1234.5
		s.step = state1
		return scanTransition_BeginLiteral
	}
	return s.error(c, "looking for beginning of value")
}

// stateBeginStringOrEmpty is the state after reading `{`.
func stateBeginStringOrEmpty(s *scanner, c byte) scanTransition {
	if c <= ' ' && isSpace(c) {
		return scanTransition_SkipSpace
	}
	if c == '}' {
		n := len(s.parseState)
		s.parseState[n-1] = parsePhase_MapValue
		return stateEndValue(s, c)
	}
	return stateBeginString(s, c)
}

// stateBeginString is the state after reading `{"key": value,`.
func stateBeginString(s *scanner, c byte) scanTransition {
	if c <= ' ' && isSpace(c) {
		return scanTransition_SkipSpace
	}
	if c == '"' {
		s.step = stateInString
		return scanTransition_BeginLiteral
	}
	return s.error(c, "looking for beginning of map key string")
}

// stateEndValue is the state after completing a value,
// such as after reading `{}` or `true` or `["x"`.
func stateEndValue(s *scanner, c byte) scanTransition {
	n := len(s.parseState)
	if n == 0 {
		// Completed top-level before the current byte.
		s.step = stateEndTop
		s.endTop = true
		return stateEndTop(s, c)
	}
	if c <= ' ' && isSpace(c) {
		s.step = stateEndValue
		return scanTransition_SkipSpace
	}
	ps := s.parseState[n-1]
	switch ps {
	case parsePhase_MapKey:
		if c == ':' {
			s.parseState[n-1] = parsePhase_MapValue
			s.step = stateBeginValue
			return scanTransition_MapKey
		}
		return s.error(c, "after map key")
	case parsePhase_MapValue:
		if c == ',' {
			s.parseState[n-1] = parsePhase_MapKey
			s.step = stateBeginString
			return scanTransition_MapValue
		}
		if c == '}' {
			s.popParseState()
			return scanTransition_EndMap
		}
		return s.error(c, "after map key:value pair")
	case parsePhase_ArrayValue:
		if c == ',' {
			s.step = stateBeginValue
			return scanTransition_ArrayValue
		}
		if c == ']' {
			s.popParseState()
			return scanTransition_EndArray
		}
		return s.error(c, "after array element")
	}
	return s.error(c, "")
}

// stateEndTop is the state after finishing the top-level value,
// such as after reading `{}` or `[1,2,3]`.
// Only space characters should be seen now.
func stateEndTop(s *scanner, c byte) scanTransition {
	if c != ' ' && c != '\t' && c != '\r' && c != '\n' {
		// Complain about non-space byte on next call.
		s.error(c, "after top-level value")
	}
	return scanTransition_End
}

// stateInString is the state after reading `"`.
func stateInString(s *scanner, c byte) scanTransition {
	if c == '"' {
		s.step = stateEndValue
		return scanTransition_Continue
	}
	if c == '\\' {
		s.step = stateInStringEsc
		return scanTransition_Continue
	}
	if c < 0x20 {
		return s.error(c, "in string literal")
	}
	return scanTransition_Continue
}

// stateInStringEsc is the state after reading `"\` during a quoted string.
func stateInStringEsc(s *scanner, c byte) scanTransition {
	switch c {
	case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
		s.step = stateInString
		return scanTransition_Continue
	case 'u':
		s.step = stateInStringEscU
		return scanTransition_Continue
	}
	return s.error(c, "in string escape code")
}

// stateInStringEscU is the state after reading `"\u` during a quoted string.
func stateInStringEscU(s *scanner, c byte) scanTransition {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU1
		return scanTransition_Continue
	}
	// numbers
	return s.error(c, "in \\u hexadecimal character escape")
}

// stateInStringEscU1 is the state after reading `"\u1` during a quoted string.
func stateInStringEscU1(s *scanner, c byte) scanTransition {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU12
		return scanTransition_Continue
	}
	// numbers
	return s.error(c, "in \\u hexadecimal character escape")
}

// stateInStringEscU12 is the state after reading `"\u12` during a quoted string.
func stateInStringEscU12(s *scanner, c byte) scanTransition {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU123
		return scanTransition_Continue
	}
	// numbers
	return s.error(c, "in \\u hexadecimal character escape")
}

// stateInStringEscU123 is the state after reading `"\u123` during a quoted string.
func stateInStringEscU123(s *scanner, c byte) scanTransition {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInString
		return scanTransition_Continue
	}
	// numbers
	return s.error(c, "in \\u hexadecimal character escape")
}

// stateNeg is the state after reading `-` during a number.
func stateNeg(s *scanner, c byte) scanTransition {
	if c == '0' {
		s.step = state0
		return scanTransition_Continue
	}
	if '1' <= c && c <= '9' {
		s.step = state1
		return scanTransition_Continue
	}
	return s.error(c, "in numeric literal")
}

// state1 is the state after reading a non-zero integer during a number,
// such as after reading `1` or `100` but not `0`.
func state1(s *scanner, c byte) scanTransition {
	if '0' <= c && c <= '9' {
		s.step = state1
		return scanTransition_Continue
	}
	return state0(s, c)
}

// state0 is the state after reading `0` during a number.
func state0(s *scanner, c byte) scanTransition {
	if c == '.' {
		s.step = stateDot
		return scanTransition_Continue
	}
	if c == 'e' || c == 'E' {
		s.step = stateE
		return scanTransition_Continue
	}
	return stateEndValue(s, c)
}

// stateDot is the state after reading the integer and decimal point in a number,
// such as after reading `1.`.
func stateDot(s *scanner, c byte) scanTransition {
	if '0' <= c && c <= '9' {
		s.step = stateDot0
		return scanTransition_Continue
	}
	return s.error(c, "after decimal point in numeric literal")
}

// stateDot0 is the state after reading the integer, decimal point, and subsequent
// digits of a number, such as after reading `3.14`.
func stateDot0(s *scanner, c byte) scanTransition {
	if '0' <= c && c <= '9' {
		return scanTransition_Continue
	}
	if c == 'e' || c == 'E' {
		s.step = stateE
		return scanTransition_Continue
	}
	return stateEndValue(s, c)
}

// stateE is the state after reading the mantissa and e in a number,
// such as after reading `314e` or `0.314e`.
func stateE(s *scanner, c byte) scanTransition {
	if c == '+' || c == '-' {
		s.step = stateESign
		return scanTransition_Continue
	}
	return stateESign(s, c)
}

// stateESign is the state after reading the mantissa, e, and sign in a number,
// such as after reading `314e-` or `0.314e+`.
func stateESign(s *scanner, c byte) scanTransition {
	if '0' <= c && c <= '9' {
		s.step = stateE0
		return scanTransition_Continue
	}
	return s.error(c, "in exponent of numeric literal")
}

// stateE0 is the state after reading the mantissa, e, optional sign,
// and at least one digit of the exponent in a number,
// such as after reading `314e-2` or `0.314e+1` or `3.14e0`.
func stateE0(s *scanner, c byte) scanTransition {
	if '0' <= c && c <= '9' {
		return scanTransition_Continue
	}
	return stateEndValue(s, c)
}

// stateT is the state after reading `t`.
func stateT(s *scanner, c byte) scanTransition {
	if c == 'r' {
		s.step = stateTr
		return scanTransition_Continue
	}
	return s.error(c, "in literal true (expecting 'r')")
}

// stateTr is the state after reading `tr`.
func stateTr(s *scanner, c byte) scanTransition {
	if c == 'u' {
		s.step = stateTru
		return scanTransition_Continue
	}
	return s.error(c, "in literal true (expecting 'u')")
}

// stateTru is the state after reading `tru`.
func stateTru(s *scanner, c byte) scanTransition {
	if c == 'e' {
		s.step = stateEndValue
		return scanTransition_Continue
	}
	return s.error(c, "in literal true (expecting 'e')")
}

// stateF is the state after reading `f`.
func stateF(s *scanner, c byte) scanTransition {
	if c == 'a' {
		s.step = stateFa
		return scanTransition_Continue
	}
	return s.error(c, "in literal false (expecting 'a')")
}

// stateFa is the state after reading `fa`.
func stateFa(s *scanner, c byte) scanTransition {
	if c == 'l' {
		s.step = stateFal
		return scanTransition_Continue
	}
	return s.error(c, "in literal false (expecting 'l')")
}

// stateFal is the state after reading `fal`.
func stateFal(s *scanner, c byte) scanTransition {
	if c == 's' {
		s.step = stateFals
		return scanTransition_Continue
	}
	return s.error(c, "in literal false (expecting 's')")
}

// stateFals is the state after reading `fals`.
func stateFals(s *scanner, c byte) scanTransition {
	if c == 'e' {
		s.step = stateEndValue
		return scanTransition_Continue
	}
	return s.error(c, "in literal false (expecting 'e')")
}

// stateN is the state after reading `n`.
func stateN(s *scanner, c byte) scanTransition {
	if c == 'u' {
		s.step = stateNu
		return scanTransition_Continue
	}
	return s.error(c, "in literal null (expecting 'u')")
}

// stateNu is the state after reading `nu`.
func stateNu(s *scanner, c byte) scanTransition {
	if c == 'l' {
		s.step = stateNul
		return scanTransition_Continue
	}
	return s.error(c, "in literal null (expecting 'l')")
}

// stateNul is the state after reading `nul`.
func stateNul(s *scanner, c byte) scanTransition {
	if c == 'l' {
		s.step = stateEndValue
		return scanTransition_Continue
	}
	return s.error(c, "in literal null (expecting 'l')")
}

// stateError is the state after reaching a syntax error,
// such as after reading `[1}` or `5.1.2`.
func stateError(s *scanner, c byte) scanTransition {
	return scanTransition_Error
}

// error records an error and switches to the error state.
func (s *scanner) error(c byte, context string) scanTransition {
	s.step = stateError
	s.err = &SyntaxError{"invalid character " + quoteChar(c) + " " + context, s.bytes}
	return scanTransition_Error
}

// quoteChar formats c as a quoted character literal
func quoteChar(c byte) string {
	// special cases - different from quoted strings
	if c == '\'' {
		return `'\''`
	}
	if c == '"' {
		return `'"'`
	}

	// use quoted string with different quotation marks
	s := strconv.Quote(string(c))
	return "'" + s[1:len(s)-1] + "'"
}

// undo causes the scanner to return scanCode from the next state transition.
// This gives callers a simple 1-byte undo mechanism.
func (s *scanner) undo(scanCode scanTransition) {
	if s.redo {
		panic("json: invalid use of scanner")
	}
	s.redoCode = scanCode
	s.redoState = s.step
	s.step = stateRedo
	s.redo = true
}

// stateRedo helps implement the scanner's 1-byte undo.
func stateRedo(s *scanner, c byte) scanTransition {
	s.redo = false
	s.step = s.redoState
	return s.redoCode
}
