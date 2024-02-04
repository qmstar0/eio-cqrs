package cqrs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio/message"
	"strings"
)

type JsonMarshaler struct {
	GenerateIDFunc func() string
}

func NewJsonMarshaler(generateIDFunc func() string) *JsonMarshaler {
	if generateIDFunc == nil {
		generateIDFunc = eio.NewUUID
	}
	return &JsonMarshaler{generateIDFunc}
}

func (j JsonMarshaler) Marshal(ctx context.Context, v any) (*message.Context, error) {
	var err error
	msg := message.WithContext(j.GenerateIDFunc(), ctx)
	msg.Payload, err = json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (j JsonMarshaler) Unmarshal(msg *message.Context, v any) error {
	return json.Unmarshal(msg.Payload, v)
}

func (j JsonMarshaler) Name(v interface{}) string {
	segments := strings.Split(fmt.Sprintf("%T", v), ".")
	return segments[len(segments)-1]
}
