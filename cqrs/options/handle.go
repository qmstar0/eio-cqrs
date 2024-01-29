package options

import (
	"github.com/qmstar0/eio-cqrs/cqrs"
)

// OnHandle
//
//	options.OnHandle(func (fn processor.HandlerFunc) processor.HandlerFunc {
//		return func(msg *message.Context) ([]*message.Context, error) {
//			//What to do before processing
//			msgs, err := fn(msg)
//			//What to do after processing
//			return msgs, err
//		}
//	})
func OnHandle(fn cqrs.HandleMessageFunc) cqrs.RouterBusOptionFunc {
	return func(bus *cqrs.RouterBus) error {
		bus.Config.HandleMessageFn = append(bus.Config.HandleMessageFn, fn)
		return nil
	}
}
