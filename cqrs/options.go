package cqrs

import (
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio/message"
	"github.com/qmstar0/eio/processor"
)

type RouterBusOptionFunc func(bus *RouterBus) error

func loadOptions(bus *RouterBus, options []RouterBusOptionFunc) error {
	for i := range options {
		err := options[i](bus)
		if err != nil {
			return err
		}
	}
	return nil
}

type GenerateTopicFunc func(string) string

func defaultGenerateTopicFn(s string) string {
	return s
}

type PublishFunc func(pub eio.Publisher, s string, msg *message.Context) error

func defaultPublishFn(pub eio.Publisher, s string, msg *message.Context) error {
	return pub.Publish(s, msg)
}

type HandleMessageFunc processor.HandlerMiddleware

func defaultHandleMessageFn(fn processor.HandlerFunc) processor.HandlerFunc {
	return func(msg *message.Context) ([]*message.Context, error) {
		return fn(msg)
	}
}
