package cqrs

import (
	"context"
	"errors"
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio/message"
	"github.com/qmstar0/eio/processor"
	"reflect"
)

type Bus interface {
	Publish(ctx context.Context, cmd any, callbacks ...Callback) error
	AddHandlers(handlers ...Handler) error
	Run(ctx context.Context) error
	Running() <-chan struct{}
}

type RouterBusConfig struct {
	OnGenerateTopic  GenerateTopicFunc
	PublishMessageFn PublishFunc

	OnHandleMessage []func(processor.HandlerFunc) processor.HandlerFunc
}

type RouterBus struct {
	Router *processor.Router

	Marshaler MessageMarshaler

	Config RouterBusConfig
}

func NewBus(publisher eio.Publisher, marshaler MessageMarshaler, opts ...RouterBusOptionFunc) (Bus, error) {
	if publisher == nil {
		return nil, errors.New("missing Publisher")
	}
	if marshaler == nil {
		return nil, errors.New("missing marshaler")
	}

	bus := &RouterBus{
		Router:    processor.NewRouter(),
		Marshaler: marshaler,
		Config: RouterBusConfig{
			OnGenerateTopic:  defaultGenerateTopicFn,
			PublishMessageFn: defaultPublishMessageFn(publisher),
			OnHandleMessage:  make([]func(processor.HandlerFunc) processor.HandlerFunc, 0),
		},
	}

	if err := loadOptions(bus, opts); err != nil {
		return nil, err
	}

	return bus, nil
}

func (c RouterBus) Run(ctx context.Context) error {
	return c.Router.Run(ctx)
}

func (c RouterBus) Running() <-chan struct{} {
	return c.Router.Running()
}

func (c RouterBus) AddHandlers(handlers ...Handler) error {
	for _, handler := range handlers {
		if err := c.addHandlerToRouter(handler); err != nil {
			return err
		}
	}

	return nil
}

func (c RouterBus) Publish(ctx context.Context, cmd any, callbacks ...Callback) error {
	var err error

	msg, err := c.newMessage(ctx, cmd)
	if err != nil {
		return err
	}

	c.setCallback(msg, callbacks)

	cmdName := c.Marshaler.Name(cmd)

	topic := c.Config.OnGenerateTopic(cmdName)

	if err = c.Config.PublishMessageFn(topic, msg); err != nil {
		return err
	}

	return nil
}

func (c RouterBus) newMessage(ctx context.Context, command any) (*message.Context, error) {
	var err error
	msg := message.WithContext(eio.NewUUID(), ctx)
	msg.Payload, err = c.Marshaler.Marshal(command)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c RouterBus) addHandlerToRouter(handler Handler) error {
	handlerName := handler.Name()
	commandName := c.Marshaler.Name(handler.SubscribedTo())
	topic := c.Config.OnGenerateTopic(commandName)

	handlerFunc, err := c.handlerToHandleFunc(handler)
	if err != nil {
		return err
	}

	for i := range c.Config.OnHandleMessage {
		handlerFunc = c.Config.OnHandleMessage[i](handlerFunc)
	}

	c.Router.AddHandler(
		handlerName,
		topic,
		handler.Subscriber(),
		handlerFunc,
	)
	return nil
}

func (RouterBus) setCallback(msg *message.Context, callbacks []Callback) {

	for _, callback := range callbacks {
		afterFn := func() {
			if errors.Is(msg.Err(), message.Done) {
				callback(msg)
			}
		}
		_ = context.AfterFunc(msg, afterFn)
	}
}

func (c RouterBus) handlerToHandleFunc(handler Handler) (processor.HandlerFunc, error) {
	to := handler.SubscribedTo()
	if !isPointerType(to) {
		return nil, SubscribedToTypeError{message: to}
	}

	return func(msg *message.Context) ([]*message.Context, error) {

		cmd := handler.SubscribedTo()

		if err := c.Marshaler.Unmarshal(msg.Payload, cmd); err != nil {
			return nil, err
		}

		err := handler.Handle(msg, cmd)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}, nil
}

func isPointerType(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}
func defaultPublishMessageFn(publisher eio.Publisher) PublishFunc {
	return func(s string, msg *message.Context) error {
		return publisher.Publish(s, msg)
	}
}

func defaultGenerateTopicFn(s string) string {
	return s
}
