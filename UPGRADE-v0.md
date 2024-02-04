<div style="text-align: center">

# eio-CQRS

</div>

## 2/4

**发现项目设计上的缺陷，准备重构😕**

### v0.2.0

**再一次重新设计了该项目**

- 以publishBus为核心的逻辑修改为以RouterBus为核心
- RouterBus依然外置
- 相比之前更方便进行多协议适配

### v0.2.1

修改了部分细节

### v0.2.2

新增中间件`WaitAndGetHandleErr()`，该方法可以获取handler在处理后返回的错误，
通过`err := bus.Publish(...)`的方式捕获handler主动返回的错误

## 2/3 v0.1.0

**已重新设计该项目**

关于之前的计划，有所修改

> - ~~考虑将`bus`核心设计为`func(topic string, msg *message.Context)`~~
> - 将内置的`Router`重新设计为外部参数

## 2/1

**正式准备重新设计该组件项目**

- 考虑将`bus`核心设计为`func(topic string, msg *message.Context)`
- 将内置的`Router`重新设计为外部参数

整个体系围绕`Router`,引用自身核心完成路由逻辑

