package options

import "github.com/qmstar0/eio-cqrs/cqrs"

func OnGenerateTopic(fn func(string) string) cqrs.RouterBusOptionFunc {
	return func(bus *cqrs.RouterBus) error {
		topicFn := bus.Config.OnGenerateTopic
		bus.Config.OnGenerateTopic = func(s string) string {
			return fn(topicFn(s))
		}
		return nil
	}
}
