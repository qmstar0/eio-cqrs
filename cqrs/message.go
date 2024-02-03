package cqrs

import (
	"context"
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio/message"
)

type MessageConstructor interface {
	EncodeMessage(ctx context.Context, data any) (*message.Context, error)
	DecodeMessage(msg *message.Context, data any) error
	Topic(data any) string
}

type MessageCodec struct {
	MessageConstructorFunc   func(ctx context.Context, data any) (*message.Context, error)
	MessageDeconstructorFunc func(msg *message.Context, data any) error
	GetMessageTopic          func(data any) string
}

func (m MessageCodec) EncodeMessage(ctx context.Context, data any) (*message.Context, error) {
	return m.MessageConstructorFunc(ctx, data)
}

func (m MessageCodec) DecodeMessage(msg *message.Context, data any) error {
	return m.MessageDeconstructorFunc(msg, data)
}

func (m MessageCodec) Topic(data any) string {
	return m.GetMessageTopic(data)
}

func NewMessageCodec(marshaler MessageMarshaler) MessageConstructor {
	msgCodec := &MessageCodec{
		MessageConstructorFunc: func(ctx context.Context, data any) (*message.Context, error) {
			var err error
			msg := message.WithContext(eio.NewUUID(), ctx)
			msg.Payload, err = marshaler.Marshal(data)
			if err != nil {
				return nil, err
			}
			return msg, nil
		},
		MessageDeconstructorFunc: func(msg *message.Context, data any) error {
			return marshaler.Unmarshal(msg.Payload, data)
		},
		GetMessageTopic: func(data any) string {
			return marshaler.Name(data)
		},
	}

	return msgCodec
}
