package server

import (
	"encoding/json"
)

func JsonAdd(data []byte, key string, value interface{}) ([]byte, error) {
	decode := make(map[string]interface{})
	decErr := json.Unmarshal(data, &decode)
	if decErr != nil {
		return data, decErr
	}
	decode[key] = value
	new_data, encErr := json.Marshal(decode)
	if encErr != nil {
		return data, encErr
	}
	return new_data, nil
}

func Encode(data Data) ([]byte, error) {
	return json.Marshal(data)
}

func Decode(data []byte) (Data, error) {
	decode := Data{}
	decErr := json.Unmarshal(data, &decode)
	if decErr != nil {
		return decode, decErr
	}
	return decode, nil
}
