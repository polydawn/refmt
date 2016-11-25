package again

import "reflect"

type VarMarshalStep func(*VarMarshalDriver, *Token) (done bool, err error)

type VarMarshalDriver struct {
	stack []VarMarshalStep
	step  VarMarshalStep
}

func stepForScan(v interface{}) VarMarshalStep {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice, reflect.Array:
		return nil // TODO
	case reflect.Map:
		switch v2 := v.(type) {
		case map[string]interface{}:
			_ = v2
			return nil // TODO special
		default:
			return nil // TODO
		}
	case reflect.Interface, reflect.Ptr:
		return nil // TODO unwrap
	case reflect.Struct:
		mach := &structScanMachine{}
		mach.init()
		return mach.Step
	case reflect.Func:
		panic("no func plz")
	default:
		panic("unreachable (kind)")
	}
}

func (vr *VarMarshalDriver) Step(tok *Token) (done bool, err error) {
	done, err = vr.step(vr, tok)
	// If the step errored: out, entirely.
	if err != nil {
		return true, err
	}
	// If the step wasn't done, return same status.
	if !done {
		return false, nil
	}
	// If it WAS done, pop next, or if stack empty, we're entirely done.
	nSteps := len(vr.stack) - 1
	if nSteps == -1 {
		return // that's all folks
	}
	vr.step = vr.stack[nSteps]
	vr.stack = vr.stack[0:nSteps]
	return false, nil
}

func (vr *VarMarshalDriver) Recurse(tok *Token, target interface{}) error {
	// Push the current stepfunc onto the stack (it's a machine entrypoint, and expected to dtrt on next call),
	// and pick a stepfunc to start in on our next item to cover.
	vr.stack = append(vr.stack, vr.step)
	vr.step = stepForScan(target)
	// Immediately make a step (we're still the delegate in charge of someone else's step).
	_, err := vr.Step(tok)
	return err
}

type mapScanMachine struct {
	obj interface{} // roughly `map[T1]T2` (no ptr needed for scan, root ref effectively is one).
	// No real choice but to have an alloc containing all keys up front,
	//  since there's no such thing as an entry iterator we can safe a ref of.
}

type structScanMachine struct {
	target interface{}
	atlas  []AtlasField // Populate on initialization.
	idx    int          // Progress marker
	value  bool         // Progress marker
}

func (sm *structScanMachine) init() {

}
func (sm *structScanMachine) Step(driver *VarMarshalDriver, tok *Token) (done bool, err error) {
	if sm.idx > len(sm.atlas) {
		panic("incorrect usage: entire struct already walked")
	}
	field := sm.atlas[sm.idx]
	if sm.value {
		driver.Recurse(tok, field.Grab(sm.target))
	} else {
		*tok = field.Name
	}
	sm.value = !sm.value
	return true, nil
}
