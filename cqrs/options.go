package cqrs

import "github.com/qmstar0/eio/message"

type CommandBusOptionFunc func(config *CommandBusConfig) error

type CommandBusOption interface {
	Option(bus *CommandBusConfig) error
}

func loadOptions(config *CommandBusConfig, options []CommandBusOptionFunc) error {
	for i := range options {
		err := options[i](config)
		if err != nil {
			return err
		}
	}
	return nil
}

type GenerateTopicFunc func(string) string

type PublishFunc func(s string, msgs ...*message.Context) error
