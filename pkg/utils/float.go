package utils

import (
	"math"
	"strconv"

	"github.com/thanhfphan/eventstore/pkg/errors"
)

func GetFloat(n interface{}) (float64, error) {
	switch i := n.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	}

	return math.NaN(), errors.New("cant convert n=%v to float64", n)
}
