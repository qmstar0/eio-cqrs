package cqrs_test

import (
	"context"
	"github.com/qmstar0/eio-cqrs/cqrs"
	"github.com/qmstar0/eio/message"
	"github.com/qmstar0/eio/processor"
	"github.com/qmstar0/eio/pubsub/gopubsub"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouter(t *testing.T) {
	ctx := getTimeoutCtx()

	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{})
	codec := cqrs.NewMessageCodec(cqrs.JsonMarshaler{})
	bus := cqrs.NewBus(pubsub, codec)

	router := processor.NewRouter()

	err := bus.WithRouter(router,
		func(fn processor.HandlerFunc) processor.HandlerFunc {
			return func(msg *message.Context) ([]*message.Context, error) {
				msg.SetValue(1, 1)
				t.Log("handler middleware before", msg)
				msgs, err := fn(msg)
				assert.Equal(t, msg.Value(1), 1)
				t.Log("handler middleware after", msg)
				return msgs, err
			}
		}).AddHandlers(cqrs.NewHandler[Cmd]("main", pubsub, func(ctx context.Context, v *Cmd) error {
		t.Log("main handler", v)
		return nil
	}))
	assert.NoError(t, err)

	go publishMessage(ctx, bus)

	err = router.Run(ctx)
	assert.NoError(t, err)
}
