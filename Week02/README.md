# error
## Error vs Exception
### Exception的特点
  * 单返值
  * 无法知道被调用方会抛出什么异常
  * checked exception
### Eorror的特点
  * 多参数返回
  * panic
  * recover

### 知识点
  * 错误的做法:来一个请求,开一个goroutine来处理
  * 正确的做法: 通过chan传递消息，交给一个工作池来统一处理
  * 因为野生的goroutine是无法被其他或主main给recover往的
  * main函数的初始化是否强依赖还是弱依赖
  * 检测配置文件时不符合预期值会panic
  * init函数初始化失败建议panic
> 只有真正意外的情况，那些表示不可恢复的程序错误，例如：索引越界、不可恢复的环境问题、栈溢出。我们才使用panic。对于其他错误情况，我们应该期望使用error来进行判定  
### Error Type
  #### Sentinel Error
  > 预定义的特定错误
  * 成为了你API的公共部分
  * 在二个包之间创建了依赖
  #### Eorro type
  > 实现了error接口的自定义类型
  * 类型断言来猎取更多上下文信息
  * 也会产生依赖
  #### panic errors
  > 不透明错误处理(Assert errors for behavior, not type)

```
type temporary interface{
    Temporary() bool   
}

func IsTemporary(err error) bool {
    te, ok := err.(temporary)
    return ok && te.Temporary()
}
```
> 这个逻辑可以在不导入定义错误的包或实际上不了解error的底层类型的情况下实现--我们只对它的行为感兴趣
### Indented flow is for errors
> 无错误的正常流程代码将成为一条直线，而不是缩进代码
## Wrap Errors
>github.com/pkg/errors
* 错误要被日志记录
* 应用程序处理错误，保证100%完整性
* 之后不再报告当前错误
### 正确使用情况
#### 在你的应用代码中，使用errors.New或者errors.Errorf返回错误
```
if err != nil {
    return err
}
```
#### 如果和其他库进行协作，考虑用errors.Wrap或者errors.Wrapf保存堆栈信息。同样适用于和标准库协作的时候
* 基础库
* 第三库
* 标准库

    直接返回错误，而不是每个错误产生的地方到处打日志
    在程序的顶部或者是工作的goroutine顶部(请求入口),使用%+v把堆栈详情记录
    使用errors.Case获取root.error，再进行和Sentinel error判定
### 总结
    wrap error是只有applications,可以选择应用的策略。具有最高可重用性的包只能返回根错误值。此机制与GO标准库中使用的相同。(kit库的sql.ErrNoRows)
    如果函数式方法不打算处理错误，那么用足够的上下文wrap errors，并将其返回到调用堆栈中
    一旦确定函数或方法将处理错误，错误不再是错误。如果函数或方法仍然需要发出返回，则它不能返回错误值，它应该只返回零    
