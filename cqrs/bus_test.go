package cqrs_test

import (
	"context"
	"github.com/qmstar0/eio-cqrs/cqrs"
	"github.com/qmstar0/eio-cqrs/cqrs/middleware"
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

	bus := cqrs.NewBus(pubsub, cqrs.NewMessageCodec(cqrs.JsonMarshaler{}))

	router := processor.NewRouter()

	routerBus := bus.WithRouter(router)
	err := routerBus.AddHandlers(cqrs.NewHandler[Cmd]("main", pubsub, func(ctx context.Context, v *Cmd) error {
		t.Log("main handler", v)
		return nil
	}))
	assert.NoError(t, err)

	go publishMessage(ctx, bus)

	err = router.Run(ctx)
	assert.NoError(t, err)
}

func TestBus_PublishOption(t *testing.T) {

	ctx := getTimeoutCtx()

	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{})

	bus := cqrs.NewBus(pubsub, cqrs.NewMessageCodec(cqrs.JsonMarshaler{}),
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

	router := processor.NewRouter()

	routerBus := bus.WithRouter(router)
	err := routerBus.AddHandlers(cqrs.NewHandler[Cmd]("main", pubsub, func(ctx context.Context, v *Cmd) error {
		t.Log("main handler", v)
		return nil
	}))
	assert.NoError(t, err)

	go publishMessage(ctx, bus)

	err = router.Run(ctx)
	assert.NoError(t, err)
}
func TestBus_PublishOption_WaitingMessageDone(t *testing.T) {

	ctx := getTimeoutCtx()

	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{})

	bus := cqrs.NewBus(pubsub, cqrs.NewMessageCodec(cqrs.JsonMarshaler{}), middleware.WaitingMessageDone())

	router := processor.NewRouter()

	routerBus := bus.WithRouter(router)
	err := routerBus.AddHandlers(cqrs.NewHandler[Cmd]("main", pubsub, func(ctx context.Context, v *Cmd) error {
		t.Log("main handler", v)
		time.Sleep(time.Second * 5)
		return nil
	}))

	assert.NoError(t, err)

	go publishMessage(ctx, bus)

	err = router.Run(ctx)
	assert.NoError(t, err)
}
