package cqrs

import (
	"context"
	"github.com/qmstar0/eio/message"
)

type MessageMarshaler interface {
	Marshal(ctx context.Context, v any) (*message.Context, error)

	Unmarshal(msg *message.Context, v any) error

	Name(v any) string
}
