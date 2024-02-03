package cqrs

import (
	"fmt"
	"github.com/qmstar0/eio/processor"
)

type RouterBus interface {
	AddHandlers(handlers ...Handler) error
}

type HandlerRouter func(handlers Handler) error

func (f HandlerRouter) AddHandlers(handlers ...Handler) error {
	var err error
	for i := range handlers {
		handler := handlers[i]
		to := handler.SubscribedTo()
		if !isPointerType(to) {
			return SubscribedToTypeError{message: to}
		}

		if err = f(handler); err != nil {
			return err
		}
	}
	return nil
}

type HandleMiddleware processor.HandlerMiddleware

func loadHandleMiddleware(handlerFunc processor.HandlerFunc, middlewares []HandleMiddleware) processor.HandlerFunc {
	handleFn := handlerFunc
	for i := range middlewares {
		handleFn = middlewares[i](handleFn)
	}
	return handleFn
}

type SubscribedToTypeError struct {
	message any
}

func (e SubscribedToTypeError) Error() string {
	return fmt.Sprintf("类型错误: %T 应该为`指针`", e.message)
}
