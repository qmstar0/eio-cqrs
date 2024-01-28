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
	Publish(ctx context.Context, cmd any) error
	WithOptions(options ...RouterBusOptionFunc) Bus
	AddHandlers(handlers ...Handler) error
	Run(ctx context.Context) error
	Running() <-chan struct{}
}

type RouterBusConfig struct {
	GenerateTopicFn GenerateTopicFunc
	PublishFn       PublishFunc

	MarshalFn   MarshalFunc
	UnMarshalFn UnMarshalFunc
}

func (c *RouterBusConfig) setDefaults() {
	if c.GenerateTopicFn == nil {
		c.GenerateTopicFn = defaultGenerateTopicFn
	}
	if c.PublishFn == nil {
		c.PublishFn = defaultPublishFn
	}
}

func (c *RouterBusConfig) Validate() error {
	var err error

	if c.GenerateTopicFn == nil {
		err = errors.Join(err, ConfigValidationError{"missing `GenerateTopicFn`"})
	}

	if c.PublishFn == nil {
		err = errors.Join(err, ConfigValidationError{"missing `PublishFn`"})
	}

	return err
}

type RouterBus struct {
	router *processor.Router

	Publisher eio.Publisher
	Marshaler MessageMarshaler

	Config RouterBusConfig
}

func NewBus(publisher eio.Publisher, marshaler MessageMarshaler, opts ...RouterBusOptionFunc) (Bus, error) {
	if publisher == nil {
		return nil, errors.New("missing publisher")
	}
	if marshaler == nil {
		return nil, errors.New("missing marshaler")
	}

	bus := &RouterBus{
		Publisher: publisher,
		router:    processor.NewRouter(),
		Marshaler: marshaler,
		Config:    RouterBusConfig{},
	}

	if err := loadOptions(bus, opts); err != nil {
		return nil, err
	}

	bus.Config.setDefaults()

	if err := bus.Config.Validate(); err != nil {
		return nil, err
	}

	return bus, nil
}

func (c RouterBus) Run(ctx context.Context) error {
	return c.router.Run(ctx)
}

func (c RouterBus) Running() <-chan struct{} {
	return c.router.Running()
}

func (c RouterBus) WithOptions(options ...RouterBusOptionFunc) Bus {
	bus := c

	if err := loadOptions(&bus, options); err != nil {
		panic(err)
	}

	bus.Config.setDefaults()

	if err := bus.Config.Validate(); err != nil {
		panic(err)
	}

	return &bus
}

func (c RouterBus) AddHandlers(handlers ...Handler) error {
	for _, handler := range handlers {
		if err := c.addHandlerToRouter(handler); err != nil {
			return err
		}
	}

	return nil
}

func (c RouterBus) Publish(ctx context.Context, cmd any) error {
	var err error

	msg, err := c.newMessage(ctx, cmd)
	if err != nil {
		return err
	}

	cmdName := c.Marshaler.Name(cmd)

	topic := c.Config.GenerateTopicFn(cmdName)

	if err = c.Config.PublishFn(c.Publisher, topic, msg); err != nil {
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
	topic := c.Config.GenerateTopicFn(commandName)

	handlerFunc, err := c.handlerToHandlerFunc(handler)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	c.router.AddHandler(
		handlerName,
		topic,
		handler.Subscriber(),
		handlerFunc,
	)
	return nil
}

func (c RouterBus) handlerToHandlerFunc(handler Handler) (processor.HandlerFunc, error) {
	to := handler.SubscribedTo()
	if reflect.ValueOf(to).Kind() != reflect.Ptr {
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
