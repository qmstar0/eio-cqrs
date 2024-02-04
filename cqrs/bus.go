package cqrs

import (
	"context"
	"fmt"
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio/message"
	"github.com/qmstar0/eio/processor"
	"reflect"
)

type PublishBus interface {
	Publish(ctx context.Context, cmd any, callbacks ...Callback) error
}
type RouterBus interface {
	AddHandlers(handlers ...Handler) error
	WithPublisher(publisher eio.Publisher, middlewares ...PublishMiddleware) PublishBus
}

type Router struct {
	AddHandlerToRouterFn func(handlers Handler) error
	marshaler            MessageMarshaler
}

func NewRouterBus(router *processor.Router, marshaler MessageMarshaler, middlewares ...HandleMiddleware) RouterBus {
	addHandlerToRouterFn := func(handler Handler) error {
		handleFn := loadHandleMiddleware(func(msg *message.Context) ([]*message.Context, error) {
			var err error
			data := handler.SubscribedTo()
			if err = marshaler.Unmarshal(msg, data); err != nil {
				return nil, err
			}

			if err = handler.Handle(msg.Context(), data); err != nil {
				return nil, err
			}
			return nil, nil
		}, middlewares)

		router.AddHandler(handler.Name(), marshaler.Name(handler.SubscribedTo()), handler.Subscriber(), handleFn)
		return nil
	}
	return &Router{
		AddHandlerToRouterFn: addHandlerToRouterFn,
		marshaler:            marshaler,
	}
}

func (r *Router) WithPublisher(publisher eio.Publisher, middlewares ...PublishMiddleware) PublishBus {
	if publisher == nil {
		panic("missing publisher")
	}
	publishMessageFn := loadPublishMiddleware(func(topic string, msg *message.Context) error {
		return publisher.Publish(topic, msg)
	}, middlewares)

	return PublishingBus(func(ctx context.Context, cmd any, callbacks []Callback) error {
		msg, err := r.marshaler.Marshal(ctx, cmd)
		if err != nil {
			return err
		}

		if len(callbacks) > 0 {
			context.AfterFunc(msg, func() { mergeCallback(callbacks...) })
		}

		return publishMessageFn(r.marshaler.Name(cmd), msg)
	})
}

func (r *Router) AddHandlers(handlers ...Handler) error {
	var err error
	for i := range handlers {
		handler := handlers[i]
		to := handler.SubscribedTo()
		if !isPointerType(to) {
			return SubscribedToTypeError{message: to}
		}

		if err = r.AddHandlerToRouterFn(handler); err != nil {
			return err
		}
	}
	return nil
}

type PublishingBus func(ctx context.Context, data any, callbacks []Callback) error

func (f PublishingBus) Publish(ctx context.Context, data any, callbacks ...Callback) error {
	return f(ctx, data, callbacks)
}

type HandleMiddleware processor.HandlerMiddleware

func loadHandleMiddleware(handlerFunc processor.HandlerFunc, middlewares []HandleMiddleware) processor.HandlerFunc {
	handleFn := handlerFunc
	for i := range middlewares {
		handleFn = middlewares[i](handleFn)
	}
	return handleFn
}

type Callback func(msg *message.Context)

func mergeCallback(callbacks ...Callback) Callback {
	return func(msg *message.Context) {
		for i := range callbacks {
			callbacks[i](msg)
		}
	}
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
func isPointerType(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}

type SubscribedToTypeError struct {
	message any
}

func (e SubscribedToTypeError) Error() string {
	return fmt.Sprintf("类型错误: %T 应该为`指针`", e.message)
}
