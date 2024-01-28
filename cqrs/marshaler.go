package cqrs

type MessageMarshaler interface {
	Marshal(v interface{}) ([]byte, error)

	Unmarshal(data []byte, v interface{}) error

	Name(v interface{}) string
}
