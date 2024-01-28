package options

import "github.com/qmstar0/eio-cqrs/cqrs"

func OnMershal(fn func(cqrs.MarshalFunc) cqrs.MarshalFunc) cqrs.RouterBusOptionFunc {
	return func(bus *cqrs.RouterBus) error {
		bus.Config.MarshalFn = fn(bus.Config.MarshalFn)
		return nil
	}
}

func OnUnMershal(fn func(cqrs.UnMarshalFunc) cqrs.UnMarshalFunc) cqrs.RouterBusOptionFunc {
	return func(bus *cqrs.RouterBus) error {
		bus.Config.UnMarshalFn = fn(bus.Config.UnMarshalFn)
		return nil
	}
}
