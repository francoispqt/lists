package maps

import "reflect"

type TesStruct struct {
	Foo string
}

func intfSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice && s.Kind() != reflect.Array {
		panic("intfSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func compare(X, Y interface{}) bool {
	xx := intfSlice(X)
	yy := intfSlice(Y)
	for k, v := range xx {
		if len(yy)-1 >= k {
			vv := yy[k]
			if vv == v {
				continue
			}
		}
		return false
	}
	return true
}
