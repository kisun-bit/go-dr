package meta

const (
	// ########### 通用错误码
	// ########### 错误码范围：0x10001 - 0x20000
	ErrInternal      = 0x10000 // 通用错误-内部标准库异常
	ErrUnknown       = 0x10001 // 通用错误-内部异常，未知错误
	ErrNoFreePort    = 0x10002 // 通用错误-无可用端口
	ErrNoFreeStorage = 0x10003 // 通用错误-无可用内存

	// ########### RS存储服务错误码
	// ########### 错误码范围：0x20001 - 0x30000
	ErrRSSStartFailed        = 0x20001 // 存储服务-启动失败
	ErrRSSConnectNatFailed   = 0x20002 // 存储服务-连接消息系统失败
	ErrRSSSpaceCollectFailed = 0x20003 // 存储服务-空间回收失败
	ErrRSSLackStorage        = 0x20004 // 存储服务-存储空间不足
	ErrRSSLackWrite          = 0x20005 // 存储服务-备份数据写入失败
	ErrRSSLackRead           = 0x20006 // 存储服务-备份数据读出失败

	// ########### RS备份服务错误码
	// ########### 错误码范围：0x30001 - 0x40000
	ErrRSBStartFailed         = 0x30001 // 备份服务-启动失败
	ErrRSBConnectNatFailed    = 0x30002 // 备份服务-连接消息系统失败
	ErrRSBConnectClientFailed = 0x30003 // 备份服务-连接客户端失败

	// ########### RS配置服务错误码
	// ########### 错误码范围：0x40001 - 0x50000

	// ########### RS消息分发服务错误码
	// ########### 错误码范围：0x50001 - 0x60000

	// ########### RS调度服务错误码
	// ########### 错误码范围：0x60001 - 0x70000

	// TODO 更多错误码类型...
)
