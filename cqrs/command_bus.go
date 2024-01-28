package cqrs

import (
	"context"
	"errors"
	"fmt"
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio/message"
	"github.com/qmstar0/eio/processor"
	"reflect"
)

type SubscriberConstructorFunc func(SubscriberConstructorParams) (eio.Subscriber, error)

type SubscriberConstructorParams struct {
	HandlerName string
	Handler     Handler
}

type CommandBusConfig struct {
	TopicPrefix           string
	GenerateTopicFn       GenerateTopicFunc
	PublishMessage        PublishFunc
	SubscriberConstructor SubscriberConstructorFunc
}

func (c *CommandBusConfig) setDefaults() {
	if c.GenerateTopicFn == nil {
		c.GenerateTopicFn = func(s string) string { return fmt.Sprintf("%s%s", c.TopicPrefix, s) }
	}
}

func (c *CommandBusConfig) Validate() error {
	var err error
	if c.SubscriberConstructor == nil {
		return ConfigValidationError{"missing SubscriberConstructor"}
	}
	return err
}

type CommandBus struct {
	//publisher eio.Publisher

	router    *processor.Router
	marshaler MessageMarshaler
	Config    CommandBusConfig
}

func NewCommandBusWithOption(marshaler MessageMarshaler, opts ...CommandBusOptionFunc) (*CommandBus, error) {
	if marshaler == nil {
		return nil, errors.New("missing marshaler")
	}

	bus := &CommandBus{
		router:    processor.NewRouter(),
		marshaler: marshaler,
		Config:    CommandBusConfig{},
	}

	if err := loadOptions(&bus.Config, opts); err != nil {
		return nil, err
	}

	bus.Config.setDefaults()

	if err := bus.Config.Validate(); err != nil {
		return nil, err
	}

	return bus, nil
}

func NewCommandBus(
	publisher eio.Publisher,
	marshaler MessageMarshaler,
	SubscriberConstructor SubscriberConstructorFunc,
) (*CommandBus, error) {
	return NewCommandBusWithOption(
		marshaler,
		func(config *CommandBusConfig) error {
			config.PublishMessage = publisher.Publish
			config.SubscriberConstructor = SubscriberConstructor
			return nil
		},
	)
}

func (c CommandBus) Run(ctx context.Context) error {
	return c.router.Run(ctx)
}

func (c CommandBus) Running() <-chan struct{} {
	return c.router.Running()
}

func (c CommandBus) WithOptions(options ...CommandBusOptionFunc) *CommandBus {
	bus := c

	if err := loadOptions(&bus.Config, options); err != nil {
		panic(err)
	}

	return &bus
}

func (c CommandBus) AddHandlers(handlers ...Handler) error {

	for _, handler := range handlers {
		if err := c.addHandlerToRouter(c.router, handler); err != nil {
			return err
		}
	}

	return nil
}

func (c CommandBus) Publish(ctx context.Context, cmd any) error {
	var err error

	msg, err := c.newMessage(ctx, cmd)
	if err != nil {
		return err
	}

	cmdName := c.marshaler.Name(cmd)

	topic := c.Config.GenerateTopicFn(cmdName)

	if err = c.Config.PublishMessage(topic, msg); err != nil {
		return err
	}

	return nil
}

func (c CommandBus) newMessage(ctx context.Context, command any) (*message.Context, error) {
	var err error
	msg := message.WithContext(eio.NewUUID(), ctx)
	msg.Payload, err = c.marshaler.Marshal(command)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c CommandBus) addHandlerToRouter(r *processor.Router, handler Handler) error {
	handlerName := handler.Name()
	commandName := c.marshaler.Name(handler.SubscribedTo())
	topic := c.Config.GenerateTopicFn(commandName)

	handlerFunc, err := c.handlerToHandlerFunc(handler)
	if err != nil {
		return err
	}

	sub, err := c.Config.SubscriberConstructor(SubscriberConstructorParams{
		HandlerName: handlerName,
		Handler:     handler,
	})
	if err != nil {
		return err
	}

	r.AddHandler(
		handlerName,
		topic,
		sub,
		handlerFunc,
	)
	return nil
}

func (c CommandBus) handlerToHandlerFunc(handler Handler) (processor.HandlerFunc, error) {
	to := handler.SubscribedTo()
	if reflect.ValueOf(to).Kind() != reflect.Ptr {
		return nil, SubscribedToTypeError{message: to}
	}
	return func(msg *message.Context) ([]*message.Context, error) {

		cmd := handler.SubscribedTo()

		if err := c.marshaler.Unmarshal(msg.Payload, cmd); err != nil {
			return nil, err
		}

		err := handler.Handle(msg, cmd)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}, nil
}
