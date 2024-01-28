package cqrs

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JsonMarshaler struct {
}

func (j JsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j JsonMarshaler) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j JsonMarshaler) Name(v interface{}) string {
	segments := strings.Split(fmt.Sprintf("%T", v), ".")
	return segments[len(segments)-1]
}
