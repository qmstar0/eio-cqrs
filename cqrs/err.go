package cqrs

import "fmt"

type DuplicateHandlerError struct {
	message string
}

func (e DuplicateHandlerError) Error() string {
	return "重复添加了handler:" + e.message
}

type ConfigValidationError struct {
	message string
}

func (e ConfigValidationError) Error() string {
	return "配置验证错误: " + e.message
}

type SubscribedToTypeError struct {
	message any
}

func (e SubscribedToTypeError) Error() string {
	return fmt.Sprintf("类型错误: %T 应该为`指针`", e.message)
}
