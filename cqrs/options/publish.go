package options

import (
	"blog/infrastructure/cqrs"
	"github.com/qmstar0/eio"
)

func OnPublish(fn func(cqrs.PublishFunc) cqrs.PublishFunc) cqrs.CommandBusOptionFunc {
	return func(config *cqrs.CommandBusConfig) error {
		config.PublishMessage = fn(config.PublishMessage)
		return nil
	}
}

func SetPublish(publisher eio.Publisher) cqrs.CommandBusOptionFunc {
	return func(config *cqrs.CommandBusConfig) error {
		config.PublishMessage = publisher.Publish
		return nil
	}
}
