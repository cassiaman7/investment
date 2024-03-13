package str

import "strconv"

func ToFloat64(s string) (num float64) {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}

	return num
}

func ToInt64(s string) (num int64) {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return num
}
