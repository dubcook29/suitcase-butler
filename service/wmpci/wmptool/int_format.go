package wmptool

import (
	"strconv"
)

func Any2Int(i interface{}) (int, bool) {
	if v, ok := i.(int); ok {
		return v, true
	} else if vv, ok := i.(int32); ok {
		return int(vv), true
	} else if vvv, ok := i.(int64); ok {
		return int(vvv), true
	}

	if v, ok := i.(float32); ok {
		return int(v), true
	} else if vv, ok := i.(float64); ok {
		return int(vv), true
	}

	if v, ok := i.(string); ok {
		if vv, err := strconv.Atoi(v); err != nil {
			return vv, true
		}
	}

	return 0, false
}
