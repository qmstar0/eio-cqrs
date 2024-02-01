<div style="text-align: center">

# eio-CQRS

</div>

## 2/1

**正式准备重新设计该组件项目**

- 考虑将`bus`核心设计为`func(topic string, msg *message.Context)`
- 将内置的`Router`重新设计为外部参数

整个体系围绕`Router`,引用自身核心完成路由逻辑

