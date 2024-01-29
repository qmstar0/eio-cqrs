package options

import "github.com/qmstar0/eio-cqrs/cqrs"

func OnGenerateTopic(fn func(string) string) cqrs.RouterBusOptionFunc {
	return func(bus *cqrs.RouterBus) error {
		if bus.Config.GenerateTopicFn == nil {
			bus.Config.GenerateTopicFn = fn
			return nil
		}

		topicFn := bus.Config.GenerateTopicFn
		bus.Config.GenerateTopicFn = func(s string) string {
			return fn(topicFn(s))
		}
		return nil
	}
}
