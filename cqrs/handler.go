package cqrs

import (
	"context"
	"github.com/qmstar0/eio"
)

type Handler interface {
	Name() string
	Subscriber() eio.Subscriber
	SubscribedTo() any
	Handle(context.Context, any) error
}

type normHandler[E any] struct {
	name     string
	sub      eio.Subscriber
	handleFn func(ctx context.Context, v *E) error
}

func NewHandler[E any](
	name string,
	subscriber eio.Subscriber,
	handleFn func(ctx context.Context, v *E) error,
) Handler {
	return &normHandler[E]{
		name:     name,
		sub:      subscriber,
		handleFn: handleFn,
	}
}

func (n normHandler[E]) Name() string {
	return n.name
}

func (n normHandler[E]) SubscribedTo() any {
	return new(E)
}

func (n normHandler[E]) Subscriber() eio.Subscriber {
	return n.sub
}

func (n normHandler[E]) Handle(ctx context.Context, v any) error {
	return n.handleFn(ctx, v.(*E))
}
