package tools

import (
	"encoding/json"
)

func StructToMap(data interface{}) (m map[string]interface{}, err error) {
	b, _ := json.Marshal(&data)
	if err = json.Unmarshal(b, &m); err != nil {
		return
	}
	return
}
