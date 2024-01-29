package cqrs_test

import (
	"context"
	"fmt"
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio-cqrs/cqrs"
	"github.com/qmstar0/eio-cqrs/cqrs/options"
	"github.com/qmstar0/eio/message"
	"github.com/qmstar0/eio/processor"
	"github.com/qmstar0/eio/pubsub/gopubsub"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Cmd struct {
	Name string
}

func getTimeoutCtx() context.Context {
	timeout, _ := context.WithTimeout(context.TODO(), time.Second*5)
	return timeout
}

func publishMessage(ctx context.Context, bus cqrs.Bus) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := bus.Publish(ctx, &Cmd{Name: "box"})
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Millisecond * 300)
		}
	}
}

func TestNewBus(t *testing.T) {

	ctx := getTimeoutCtx()

	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{})

	bus, err := cqrs.NewBus(pubsub, cqrs.JsonMarshaler{})
	assert.NoError(t, err)

	err = bus.AddHandlers(cqrs.NewHandler[Cmd]("handler", pubsub, func(ctx context.Context, v *Cmd) error {
		t.Logf("%s", v)
		return nil
	}))
	assert.NoError(t, err)

	go publishMessage(ctx, bus)

	err = bus.Run(ctx)
	assert.NoError(t, err)
}

func TestRouterBus_WithOptions(t *testing.T) {
	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{
		MessageChannelBuffer: 5,
	})

	bus, err := cqrs.NewBus(pubsub, cqrs.JsonMarshaler{})
	assert.NoError(t, err)

	err = bus.AddHandlers(cqrs.NewHandler[Cmd]("handler1", pubsub, func(ctx context.Context, v *Cmd) error {
		fmt.Println("handler1 handle msg:", v)
		return nil
	}))
	assert.NoError(t, err)

	t.Run("OnGenerateTopic", func(t *testing.T) {
		ctx := getTimeoutCtx()
		newbus := bus.WithOptions(
			options.OnGenerateTopic(func(s string) string {
				return "test_" + s
			}),
		)

		err = newbus.AddHandlers(cqrs.NewHandler[Cmd]("handler2", pubsub, func(ctx context.Context, v *Cmd) error {
			t.Logf("handler2 handle cmd: %s", v)
			return nil
		}))
		assert.NoError(t, err)

		go publishMessage(ctx, newbus)
		err = newbus.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("OnPublish", func(t *testing.T) {
		ctx := getTimeoutCtx()
		newbus := bus.WithOptions(
			options.OnPublish(func(publishFunc cqrs.PublishFunc) cqrs.PublishFunc {
				return func(pub eio.Publisher, s string, msg *message.Context) error {
					t.Log("payload", msg.Payload)
					return publishFunc(pub, s, msg)
				}
			}))

		go publishMessage(ctx, newbus)
		err := newbus.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("OnPublish+OnGenerateTopic", func(t *testing.T) {
		ctx := getTimeoutCtx()
		newbus := bus.WithOptions(
			options.OnPublish(func(publishFunc cqrs.PublishFunc) cqrs.PublishFunc {
				return func(pub eio.Publisher, s string, msg *message.Context) error {
					t.Log("payload", msg.Payload)
					return publishFunc(pub, s, msg)
				}
			}),
			options.OnGenerateTopic(func(s string) string {
				return "test_" + s
			}),
		)

		go publishMessage(ctx, newbus)
		err := newbus.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("OnHandle", func(t *testing.T) {
		ctx := getTimeoutCtx()
		newbus := bus.WithOptions(
			options.OnHandle(func(fn processor.HandlerFunc) processor.HandlerFunc {
				return func(msg *message.Context) ([]*message.Context, error) {
					t.Log("OnHandle msg:", msg)
					return fn(msg)
				}
			}),
		)

		err = newbus.AddHandlers(cqrs.NewHandler[Cmd]("handler3", pubsub, func(ctx context.Context, v *Cmd) error {
			t.Logf("handler3 handle cmd: %s", v)
			return nil
		}))

		go publishMessage(ctx, newbus)
		err := newbus.Run(ctx)
		assert.NoError(t, err)
	})

}
