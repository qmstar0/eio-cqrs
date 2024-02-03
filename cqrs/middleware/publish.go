package middleware

import (
	"github.com/qmstar0/eio-cqrs/cqrs"
	"github.com/qmstar0/eio/message"
)

func WaitingMessageDone() cqrs.PublishMiddleware {
	return func(publishFunc cqrs.PublishFunc) cqrs.PublishFunc {
		return func(topic string, msg *message.Context) error {
			err := publishFunc(topic, msg)
			if err != nil {
				return err
			}

			<-msg.Done()
			return nil
		}
	}
}