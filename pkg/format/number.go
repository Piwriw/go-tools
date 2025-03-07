package format

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		if v > math.MaxFloat64 {
			return 0, fmt.Errorf("value too large to convert: %d", v)
		}
		return float64(v), nil
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return 0, fmt.Errorf("cannot convert json.Number to float64: %v", err)
		}
		return f, nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string to float64: %v", err)
		}
		return f, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, fmt.Errorf("cannot convert nil to float64")
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}
