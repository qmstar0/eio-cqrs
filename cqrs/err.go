package cqrs

import "fmt"

type SubscribedToTypeError struct {
	message any
}

func (e SubscribedToTypeError) Error() string {
	return fmt.Sprintf("类型错误: %TC 应该为`指针`", e.message)
}
