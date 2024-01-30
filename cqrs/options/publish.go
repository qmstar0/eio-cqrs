package options

import (
	"github.com/qmstar0/eio-cqrs/cqrs"
)

// OnPublish
//
//	options.OnPublish(func (publishFunc cqrs.PublishFunc) cqrs.PublishFunc {
//		return func(s string, msg *message.Context) error {
//			t.Log("before publish log msg:", msg)
//			err := publishFunc(s, msg)
//			t.Log("before publish log msg:", msg)
//			return err
//		}
//	})
func OnPublish(fn func(cqrs.PublishFunc) cqrs.PublishFunc) cqrs.RouterBusOptionFunc {
	return func(bus *cqrs.RouterBus) error {
		bus.PublishMessageFn = fn(bus.PublishMessageFn)
		return nil
	}
}
