package cqrs_test

import (
	"context"
	"github.com/qmstar0/eio-cqrs/cqrs"
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

	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{
		MessageChannelBuffer: 5,
	})

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
