package cqrs_test

import (
	"context"
	"github.com/qmstar0/eio-cqrs/cqrs"
	"github.com/qmstar0/eio-cqrs/cqrs/options"
	"github.com/qmstar0/eio/message"
	"github.com/qmstar0/eio/processor"
	"github.com/qmstar0/eio/pubsub/gopubsub"
	"github.com/stretchr/testify/assert"
	"strings"
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

	bus, err := cqrs.NewBus(pubsub, cqrs.JsonMarshaler{},
		options.OnGenerateTopic(func(s string) string {
			return "test1_" + s
		}),
		options.OnGenerateTopic(func(s string) string {
			return "test2_" + s
		}),
		options.OnPublish(func(publishFunc cqrs.PublishFunc) cqrs.PublishFunc {
			return func(s string, msg *message.Context) error {
				assert.True(t, strings.HasPrefix(s, "test2_test1_"))
				//t.Log("before publish log msg:", s, msg)
				msg.SetValue(1, 1)
				err := publishFunc(s, msg)
				//t.Log("before publish log msg:", s, msg)
				return err
			}
		}),
		options.OnHandle(func(fn processor.HandlerFunc) processor.HandlerFunc {
			return func(msg *message.Context) ([]*message.Context, error) {
				//t.Log("before processing log msg:", msg)
				assert.Equal(t, msg.Value(1), 1)
				msgs, err := fn(msg)
				//t.Log("after processing log msg:", msg)
				assert.Equal(t, msg.Value(1), 1)
				return msgs, err
			}
		}),
	)
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
