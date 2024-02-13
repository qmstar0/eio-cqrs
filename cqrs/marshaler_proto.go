package cqrs

import (
	"context"
	"errors"
	"fmt"
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio/message"
	"google.golang.org/protobuf/proto"
	"strings"
)

type protoMarshaler struct {
	GenerateIDFunc func() string
}

func NewProtoMarshaler(GenerateIDFunc func() string) MessageMarshaler {
	if GenerateIDFunc == nil {
		GenerateIDFunc = eio.NewUUID
	}
	return protoMarshaler{GenerateIDFunc: GenerateIDFunc}
}

func (p protoMarshaler) Marshal(ctx context.Context, v any) (*message.Context, error) {
	var err error
	msg := message.WithContext(p.GenerateIDFunc(), ctx)
	pmsg, ok := v.(proto.Message)
	if !ok {
		return nil, errors.New("you have to implement `proto.Message` interface")
	}
	msg.Payload, err = proto.Marshal(pmsg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (p protoMarshaler) Unmarshal(msg *message.Context, v any) error {
	return proto.Unmarshal(msg.Payload, v.(proto.Message))
}

func (p protoMarshaler) Name(v any) string {
	segments := strings.Split(fmt.Sprintf("%T", v), ".")
	return segments[len(segments)-1]
}
