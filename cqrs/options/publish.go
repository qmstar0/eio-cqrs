package options

import (
	"github.com/qmstar0/eio"
	"github.com/qmstar0/eio-cqrs/cqrs"
)

func OnPublish(fn func(cqrs.PublishFunc) cqrs.PublishFunc) cqrs.RouterBusOptionFunc {
	return func(bus *cqrs.RouterBus) error {
		bus.Config.PublishFn = fn(bus.Config.PublishFn)
		return nil
	}
}

func SetPublish(publisher eio.Publisher) cqrs.RouterBusOptionFunc {
	return func(bus *cqrs.RouterBus) error {
		bus.Publisher = publisher
		return nil
	}
}
