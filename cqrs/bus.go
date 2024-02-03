package cqrs

import (
	"context"
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio/message"
	"github.com/qmstar0/eio/processor"
	"reflect"
)

type Callback func(msg *message.Context)

type Bus interface {
	Publish(ctx context.Context, cmd any, callbacks ...Callback) error
	WithRouter(router *processor.Router, middlewares ...HandleMiddleware) RouterBus
}

type PublishingBus struct {
	publishFn   func(ctx context.Context, data any, callbacks []Callback) error
	constructor MessageConstructor
}

func (p PublishingBus) Publish(ctx context.Context, data any, callbacks ...Callback) error {
	return p.publishFn(ctx, data, callbacks)
}

func (p PublishingBus) WithRouter(router *processor.Router, middlewares ...HandleMiddleware) RouterBus {
	return HandlerRouter(func(handler Handler) error {
		handleFn := loadHandleMiddleware(func(msg *message.Context) ([]*message.Context, error) {
			var err error
			data := handler.SubscribedTo()
			if err = p.constructor.DecodeMessage(msg, data); err != nil {
				return nil, err
			}

			if err = handler.Handle(msg.Context(), data); err != nil {
				return nil, err
			}
			return nil, nil
		}, middlewares)

		router.AddHandler(handler.Name(), p.constructor.Topic(handler.SubscribedTo()), handler.Subscriber(), handleFn)
		return nil
	})
}

func NewBus(publisher eio.Publisher, constructor MessageConstructor, middlewares ...PublishMiddleware) Bus {
	if publisher == nil {
		panic("missing publisher")
	}
	if constructor == nil {
		panic("missing constructor")
	}

	publishMessageFn := loadPublishMiddleware(func(topic string, msg *message.Context) error {
		return publisher.Publish(topic, msg)
	}, middlewares)

	return &PublishingBus{
		publishFn: func(ctx context.Context, cmd any, callbacks []Callback) error {
			msg, err := constructor.EncodeMessage(ctx, cmd)
			if err != nil {
				return err
			}

			if len(callbacks) > 0 {
				context.AfterFunc(msg, func() { mergeCallback(callbacks...) })
			}

			return publishMessageFn(constructor.Topic(cmd), msg)
		},
		constructor: constructor,
	}
}

func mergeCallback(callbacks ...Callback) Callback {
	return func(msg *message.Context) {
		for i := range callbacks {
			callbacks[i](msg)
		}
	}
}

func isPointerType(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}

type PublishFunc func(topic string, msg *message.Context) error

type PublishMiddleware func(PublishFunc) PublishFunc

func loadPublishMiddleware(publishFunc PublishFunc, middlewares []PublishMiddleware) PublishFunc {
	publishFn := publishFunc
	for i := range middlewares {
		publishFn = middlewares[i](publishFn)
	}
	return publishFn
}
