package options

import (
	"github.com/qmstar0/eio-cqrs/cqrs"
)

// OnHandle
//
//	options.OnHandle(func (fn processor.HandlerFunc) processor.HandlerFunc {
//		return func(msg *message.Context) ([]*message.Context, error) {
//			t.Log("before processing log msg:", msg)
//			msgs, err := fn(msg)
//			t.Log("after processing log msg:", msg)
//			return msgs, err
//		}
//	})
func OnHandle(fn cqrs.HandleMessageFunc) cqrs.RouterBusOptionFunc {
	return func(bus *cqrs.RouterBus) error {
		bus.Config.OnHandleMessage = append(bus.Config.OnHandleMessage, fn)
		return nil
	}
}
