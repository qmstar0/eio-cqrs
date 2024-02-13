package cqrs_test

import (
	"context"
	"github.com/qmstar0/eio-cqrs/cqrs"
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
	timeout, _ := context.WithTimeout(context.TODO(), time.Second*3)
	return timeout
}

func publishMessage(ctx context.Context, bus cqrs.PublishBus) {
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

func TestNewRouterBus(t *testing.T) {

	ctx := getTimeoutCtx()

	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{})

	router := processor.NewRouter()

	routerBus := cqrs.NewRouterBus(router, cqrs.NewJsonMarshaler(nil))

	bus := routerBus.WithPublisher(pubsub)

	err := routerBus.AddHandlers(cqrs.NewHandler[Cmd]("main", pubsub, func(ctx context.Context, v Cmd) error {
		t.Log("main handler", v)
		return nil
	}))
	assert.NoError(t, err)

	go publishMessage(ctx, bus)

	err = router.Run(ctx)
	assert.NoError(t, err)
}

func TestPublishBusMiddleware(t *testing.T) {

	ctx := getTimeoutCtx()

	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{})

	router := processor.NewRouter()

	routerBus := cqrs.NewRouterBus(router, cqrs.NewJsonMarshaler(nil))
	err := routerBus.AddHandlers(cqrs.NewHandler[Cmd]("main", pubsub, func(ctx context.Context, v Cmd) error {
		t.Log("main handler", v)
		return nil
	}))
	assert.NoError(t, err)

	bus := routerBus.WithPublisher(pubsub,
		func(publishFunc cqrs.PublishFunc) cqrs.PublishFunc {
			return func(topic string, msg *message.Context) error {
				t.Log("publish1 option before", topic, msg)
				err := publishFunc(topic, msg)
				t.Log("publish1 option after", topic, msg)
				return err
			}
		},
		func(publishFunc cqrs.PublishFunc) cqrs.PublishFunc {
			return func(topic string, msg *message.Context) error {
				t.Log("publish2 option before", topic, msg)
				err := publishFunc(topic, msg)
				t.Log("publish2 option after", topic, msg)
				return err
			}
		},
	)

	go publishMessage(ctx, bus)

	err = router.Run(ctx)
	assert.NoError(t, err)
}
func TestRouterBusMiddleware(t *testing.T) {
	ctx := getTimeoutCtx()

	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{})

	router := processor.NewRouter()

	routerBus := cqrs.NewRouterBus(router, cqrs.NewJsonMarshaler(nil),
		func(fn processor.HandlerFunc) processor.HandlerFunc {
			return func(msg *message.Context) ([]*message.Context, error) {
				msg.SetValue(1, 1)
				t.Log("handler middleware before", msg)
				msgs, err := fn(msg)
				assert.Equal(t, msg.Value(1), 1)
				t.Log("handler middleware after", msg)
				return msgs, err
			}
		})
	err := routerBus.AddHandlers(cqrs.NewHandler[Cmd]("main", pubsub, func(ctx context.Context, v Cmd) error {
		t.Log("main handler", v)
		return nil
	}))
	assert.NoError(t, err)

	bus := routerBus.WithPublisher(pubsub)

	go publishMessage(ctx, bus)

	err = router.Run(ctx)
	assert.NoError(t, err)
}
