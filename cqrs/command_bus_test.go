package cqrs_test

import (
	"context"
	"eio-cqrs/cqrs"
	"eio-cqrs/cqrs/options"
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio/message"
	"github.com/qmstar0/eio/pubsub/gopubsub"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Cmd struct {
	Name string
}

func TestNewCommandBus(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{})

	bus, err := cqrs.NewCommandBus(pubsub, cqrs.JsonMarshaler{}, func(params cqrs.SubscriberConstructorParams) (eio.Subscriber, error) {
		return pubsub, nil
	})
	assert.NoError(t, err)

	err = bus.AddHandlers(cqrs.NewHandler[Cmd]("test", func(ctx context.Context, v *Cmd) error {
		t.Log("test")
		return nil
	}))
	assert.NoError(t, err)

	go func() {
		<-bus.Running()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				err = bus.Publish(ctx, &Cmd{Name: "bob"})
				assert.NoError(t, err)
				time.Sleep(time.Millisecond * 300)
			}
		}

	}()

	err = bus.Run(ctx)
	assert.NoError(t, err)
}

func TestNewCommandBusWithOption(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pubsub := gopubsub.NewGoPubsub("test", gopubsub.GoPubsubConfig{})

	b, err := cqrs.NewCommandBusWithOption(cqrs.JsonMarshaler{}, options.SetPublish(pubsub), options.SetSubscriber(pubsub))
	assert.NoError(t, err)

	bus := b.WithOptions(options.OnPublish(func(publishFunc cqrs.PublishFunc) cqrs.PublishFunc {
		return func(s string, msgs ...*message.Context) error {
			t.Logf("%s\n", s)
			return publishFunc(s, msgs...)
		}
	}))

	err = bus.AddHandlers(cqrs.NewHandler[Cmd]("test", func(ctx context.Context, v *Cmd) error {
		t.Log("test")
		return nil
	}))
	assert.NoError(t, err)

	go func() {
		<-bus.Running()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				err = bus.Publish(ctx, &Cmd{Name: "bob"})
				assert.NoError(t, err)
				time.Sleep(time.Millisecond * 300)
			}
		}

	}()

	err = bus.Run(ctx)
	assert.NoError(t, err)
}
