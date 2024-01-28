package options

import (
	"blog/infrastructure/cqrs"
	"github.com/qmstar0/eio"
)

func SetSubscriber(subscriber eio.Subscriber) cqrs.CommandBusOptionFunc {
	return func(config *cqrs.CommandBusConfig) error {
		config.SubscriberConstructor = func(params cqrs.SubscriberConstructorParams) (eio.Subscriber, error) {
			return subscriber, nil
		}
		return nil
	}
}

func SetSubscriberConstructor(constructorFunc cqrs.SubscriberConstructorFunc) cqrs.CommandBusOptionFunc {
	return func(config *cqrs.CommandBusConfig) error {
		config.SubscriberConstructor = constructorFunc
		return nil
	}
}
