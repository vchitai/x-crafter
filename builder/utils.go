package builder

import "encoding/json"

func cloneMap(m map[string]interface{}) map[string]interface{} {
	return interfaceToMap(m)
}

func interfaceToMap(i interface{}) map[string]interface{} {
	b, err := json.Marshal(i)
	if err != nil {
		return nil
	}
	var res map[string]interface{}
	if err := json.Unmarshal(b, &res); err != nil {
		return nil
	}
	return res
}
func interfaceToMapSlice(i interface{}) []map[string]interface{} {
	b, err := json.Marshal(i)
	if err != nil {
		return nil
	}
	var res []map[string]interface{}
	if err := json.Unmarshal(b, &res); err != nil {
		return nil
	}
	return res
}
