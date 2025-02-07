package utils

import (
	"bytes"
	"encoding/json"
)

// convert a map to struct like
func FillStruct(s interface{}, mp map[string]interface{}) error {
	buf := new(bytes.Buffer)

	// encode data to struct
	err := json.NewEncoder(buf).Encode(mp)

	if err != nil {
		return err
	}
	// decode data to struct
	err = json.NewDecoder(buf).Decode(s)
	return err
}
