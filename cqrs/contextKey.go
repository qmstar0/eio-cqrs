package cqrs

var publisherBusName = &contextKey{"publisherBusName"}

type contextKey struct {
	message string
}

func (k *contextKey) String() string {
	return "context value " + k.message
}
