package cqrs

import (
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

type PublishFunc func(s string, msg *message.Context) error

type HandleMessageFunc processor.HandlerMiddleware

type Callback func(msg *message.Context)
