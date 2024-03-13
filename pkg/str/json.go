package str

import "encoding/json"

func ToJSON(data interface{}) string {
	dataByte, _ := json.MarshalIndent(data, "", "  ")
	return string(dataByte)
}
