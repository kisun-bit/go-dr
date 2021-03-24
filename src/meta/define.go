/* 负责公共库的下属项目的统一配置定义
ex: 日志配置、运行时环境配置、服务配置.....
*/

package meta

const (
	// ######## 日志配置
	LOGRootDirInLinux = "/var/log/rs"             // linux下日志文件存储的根目录
	LOGRootDirInWin   = `D:\Program Files\rs\log` // windows下日志文件存储的根目录

	// ######## 运行时环境配置
	ErrDefaultFrameLevel = 10
)
