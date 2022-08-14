package service

import (
	"bytes"
	"encoding/json"
)

func Dump(i interface{}) (s string) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetIndent(" ", "  ")
	if err := enc.Encode(i); err != nil {
		return "null"
	}
	return buf.String()
}
