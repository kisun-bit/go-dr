# `Jpkt公共库`-开发文档（Go）

`Jpkt公共库`为精容数安企业公共库代码，用于存放公共逻辑代码，在一定程度上减少“重复代码”，同时引领编码风格的规范，该代码库由大家共同维护

# 介绍篇

## 一、结构

```
├─demo <-------------------公共库测试代码，入库的代码必须保证通过测试
│  ├─core
│  └─log
├─docs <-------------------公共库的使用文档及说明
└─src
    ├─core <---------------存放核心公共代码
        ├──debug.go <-----------debug.go     运行时调试
        ├──encrypt.go <---------encrypt.go   加密
        ├──exception.go <-------exception.go “标准错误”
        ├──lock.go <------------lock.go      "锁"，不限于自旋、文件锁(跨进程使用)等
        ├──osaccess.go <--------osaccess.go  系统相关，命令执行、命令解析、状态查询等等
        ├──pool <---------------pool.go      任务池相关
    ├─datahandle
        ├──fmt.go <-------------fmt.go       负责格式化数据的公共通用代码
        ├──parse.go <-----------parse.go     数据结构解析工具
    ├─grpc
        ├── TODO // ？？？？后续业务扩充，涉及多服务化，这里提供一种简单创建gRPC C/S的方式
    ├─log
        ├──log.go <-------------log.go       统一日志工具，使用方法下文介绍
    └─meta
        ├──define.go <----------define.go    元数据类型，统一常量配置定义
        ├──errcode.go <---------errcode.go   标准化错误码
        
......
```

！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！
**TODO**： 上述是首版结构，后续可对该结构进行拆分、重构、自定制或补充



# 使用篇

## 一、日志库的使用

### 背景

统一日志使用规范，该日志库集成了`zap`高性能日志库，在性能上表现强劲，支持日志文件的切割、归档、按级别存储等功能；

![image-20210323155349321](C:\Users\kisun\AppData\Roaming\Typora\typora-user-images\image-20210323155349321.png)

### 使用方法

**REMARK**: 暂时使用导入本地包的方式来展开说明，后续通过安装至`GOPATH`中，直接引用该包

这里以举例子方式展开说明：

**第一步：**存在下述结构的项目目录，通过`go.mod`导入`Jpkt`

![image-20210323152924193](C:\Users\kisun\AppData\Roaming\Typora\typora-user-images\image-20210323152924193.png)

```mod
# go.mod 内容
module JLogDemo

go 1.14

require "jpkt" v0.0.0
replace "jpkt" => "../Jpkt"
```

**第二步**：编写日志配置文件`logging.toml`，可自己选择放于任意目录下，只要目录的可读性强就可，  内容如下:

```toml
[default]    # logger名，不能重复
Filename      = "./jlogdemo.log"  # 日志文件的存储位置
MaxSize       = 30                # 单个日志文件的最大容量，单位：MB
MaxBackups    = 6                 # 最大保留多少个日志归档文件
MaxAge        = 30                # 最大保留天数，单位：天数
Compress      = false             # 是否开启日志压缩，为了可读性一般使用false
Level         = "debug"           # 日志文件打印的最小日志级别，低于该级别便不打印
SplitByLevel  = false             # 是否按日志级别分开存储
```

**第三步：**初始化日志"对象"`Logger` - 编写`logging`目录下的`logging.go`

```go
package logging

import (
	"go.uber.org/zap"
	"jpkt/src/log"
)

var (
	Logger *zap.Logger
)

func init() {

	// Logger1
	Logger = log.GetJLoggerByConf(`D:\workspace\jrsa\JLogDemo`, "logging", "default")
	defer Logger.Sync()

	// Logger2
	// TODO

	// Logger3
	// TODO
}
```

**第四步**：使用时导入包，并使用`Logger`即可：

![image-20210323160056461](C:\Users\kisun\AppData\Roaming\Typora\typora-user-images\image-20210323160056461.png)

**第五步：**查看日志输出（`./jlogdemo.log`）

![image-20210323165310378](C:\Users\kisun\AppData\Roaming\Typora\typora-user-images\image-20210323165310378.png)



## 二、标准错误及异常

### 背景

统一错误及异常的定义、处理等

### 使用方法

**已完成的错误定义：**

```go
const (
	// ########### 通用错误码
	// ########### 错误码范围：0x10001 - 0x20000
	JErrInternal      = 0x10000 // 通用错误-内部标准库异常
	JErrUnknown       = 0x10001 // 通用错误-内部异常，未知错误
	JErrNoFreePort    = 0x10002 // 通用错误-无可用端口
	JErrNoFreeStorage = 0x10003 // 通用错误-无可用内存

	// ########### RS存储服务错误码
	// ########### 错误码范围：0x20001 - 0x30000
	JErrRSSStartFailed        = 0x20001 // 存储服务-启动失败
	JErrRSSConnectNatFailed   = 0x20002 // 存储服务-连接消息系统失败
	JErrRSSSpaceCollectFailed = 0x20003 // 存储服务-空间回收失败
	JErrRSSLackStorage        = 0x20004 // 存储服务-存储空间不足
	JErrRSSLackWrite          = 0x20005 // 存储服务-备份数据写入失败
	JErrRSSLackRead           = 0x20006 // 存储服务-备份数据读出失败

	// ########### RS备份服务错误码
	// ########### 错误码范围：0x30001 - 0x40000
	JErrRSBStartFailed         = 0x30001 // 备份服务-启动失败
	JErrRSBConnectNatFailed    = 0x30002 // 备份服务-连接消息系统失败
	JErrRSBConnectClientFailed = 0x30003 // 备份服务-连接客户端失败

	// ########### RS配置服务错误码
	// ########### 错误码范围：0x40001 - 0x50000

	// ########### RS消息分发服务错误码
	// ########### 错误码范围：0x50001 - 0x60000

	// ########### RS调度服务错误码
	// ########### 错误码范围：0x60001 - 0x70000

	// TODO ... 添加更多错误定义
)
```

**已完成的错误及异常处理：**

* ```go
  JpktStandardError
  ```

  【通用标准错误类型】

  包括`错误码`、`错误类型`、`错误原因(用户可读)`、`错误调试(工程师可读)`、`动态调用栈`等五组属性，其中错误码需严格使用公共库中的错误码定义，动态调用栈的最大长度为10

  * ```go
    func (jse *JpktStandardError) ErrorDetail() string
    // 格式化输出错误的详细信息：错误码、错误信息、错误调试、调用栈（最大10层）
    
    // 示例ex:
    //    ----------
    //    Traceback (most recent call last):
    //        [Tance] 0x10000|InternalError
    //        File "<D:/workspace/jrsa/Jpkt/src/core/exception.go>", Line 42 ----> //jpkt/src/core.(*JpktStandardError).createErrCtx
    //        File "<D:/workspace/jrsa/Jpkt/src/core/exception.go>", Line 131 ----> //jpkt/src/core.RaiseStandardError
    //        File "<D:/workspace/jrsa/Jpkt/src/core/exception.go>", Line 139 ----> //jpkt/src/core.StandardizeErr
    //        File "<D:/workspace/jrsa/Jpkt/demo/log/demo_log.go>", Line 55 ----> main.test4
    //        File "<D:/workspace/jrsa/Jpkt/demo/log/demo_log.go>", Line 79 ----> main.main
    //        File "<D:/env/go/1.14.4/src/runtime/proc.go>", Line 203 ----> runtime.main
    //        File "<D:/env/go/1.14.4/src/runtime/asm_amd64.s>", Line 1373 ----> //runtime.goexit
    //        ......
    //    ####Code: "0x10000", Type: "InternalError"
    //    ####Debug: "EEEEEEEEEEEEEEE", Desc: "内部错误，错误代码：0x10000"
    //    ----------
    ```

  * ```go
    func (jse *JpktStandardError) Error() string
    // 错误信息的简短描述，格式：仅包含ErrorDebug及ErrorCode
    
    // 示例ex:
    //     ！！！！！JpktStandardError ---> ErrCode: "0x10002", 	ErrDebug: "debug"
    ```

* ```go
  RaiseStandardError(code uint32, errType, reason, debug string) (err *JpktStandardError)
  ```

  【生成标准错误】

* ```go
  StandardizeErr(err error) (jse *JpktStandardError)
  ```

  【将普通错误类型转换为标准错误】

* ```go
  CatchPanicErr(logger *zap.Logger) *JpktStandardError
  ```

  【捕获Panic异常，并返回标准异常，常与defer联用，用于“吃掉”异常】

* ```go
  StandardPanic(code uint32, errType, reason, debug string)
  ```

  【抛出标准异常并结束程序执行】

更多使用请查阅`./demo/core/err.log`