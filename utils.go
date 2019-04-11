package eltee

import (
	"fmt"
)

func ValAsFloat(v interface{}) float64 {
	r, ok := v.(float64)
	if !ok {
		r, ok := v.(float32)
		if !ok {
			// BUT, it might be an int which we could cast as a float?
			// So use the result from valAsInt to try to get a numerical
			// value out of this thing
			return float64(ValAsInt(v))
		}
		return float64(r)
	}
	return r
}

func ValAsString(v interface{}) string {
	r, ok := v.(string)
	if !ok {
		r, ok := v.(fmt.Stringer)
		if !ok {
			return ""
		}
		return r.String()
	}
	return r
}

func ValAsInt(v interface{}) int {
	r, ok := v.(int64)
	if !ok {
		r, ok := v.(int32)
		if !ok {
			r, ok := v.(int)
			if !ok {
				return 0
			}
			return r
		}
		return int(r)
	}
	return int(r)
}
