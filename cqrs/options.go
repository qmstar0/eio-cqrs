package cqrs

import (
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio/message"
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

type PublishFunc func(pub eio.Publisher, s string, msgs ...*message.Context) error

func defaultPublishFn(pub eio.Publisher, s string, msgs ...*message.Context) error {
	return pub.Publish(s, msgs...)
}

type MarshalFunc func(marshaler MessageMarshaler, v any) ([]byte, error)
type UnMarshalFunc func(marshaler MessageMarshaler, data []byte, v any) error

func defaultMarshalFn(marshaler MessageMarshaler, v any) ([]byte, error) {
	return marshaler.Marshal(v)
}

func defaultUnMarshalFn(marshaler MessageMarshaler, data []byte, v any) error {
	return marshaler.Unmarshal(data, v)
}
