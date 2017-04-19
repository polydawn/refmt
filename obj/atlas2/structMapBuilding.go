package atlas

import "reflect"

func (sm *StructMap) AddField(fieldName string, serialName string) {

}

func fieldNameToReflectRoute(rt reflect.Type, fieldName []string) (rr reflectRoute, err error) {
	for _, fn := range fieldName {
		rf, ok := rt.FieldByName(fn)
		if !ok {
			return nil, ErrStructureMismatch{rt.Name(), "does not have field named " + fn}
		}
		rr = append(rr, rf.Index...)
	}
	return rr, nil
}
