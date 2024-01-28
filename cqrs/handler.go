package cqrs

import (
	"context"
)

type Handler interface {
	Name() string
	SubscribedTo() any
	Handle(context.Context, any) error
}

type normHandler[E any] struct {
	name     string
	handleFn func(ctx context.Context, v *E) error
}

func NewHandler[E any](
	name string,
	handleFn func(ctx context.Context, v *E) error,
) Handler {
	return &normHandler[E]{
		name:     name,
		handleFn: handleFn,
	}
}

func (n normHandler[E]) Name() string {
	return n.name
}

func (n normHandler[E]) SubscribedTo() any {
	return new(E)
}

func (n normHandler[E]) Handle(ctx context.Context, v any) error {
	return n.handleFn(ctx, v.(*E))
}
