package str

import "time"

func ToDate(t time.Time) string {
	return t.Format("2006-01-02")
}
