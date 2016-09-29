package xlate

import (
	"fmt"
)

func capturePanics(fn func()) (e error) {
	defer func() {
		if rcvr := recover(); rcvr != nil {
			e = rcvr.(error)
		}
	}()
	fn()
	return
}

func stringyEquality(x, y interface{}) bool {
	return fmt.Sprintf("%#v", x) == fmt.Sprintf("%#v", y)
}
