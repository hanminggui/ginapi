package ginapi

import (
	"errors"
	. "github.com/hanminggui/ginapi/common"
	. "github.com/hanminggui/ginapi/log"
)

// Start 入口
func Start(configPath string) {
	if len(apiStucts) < 1 {
		FatalError(errors.New("未加载任何API，请先调用LoadApi"))
	}
	// 加载配置
	initConfig(configPath)
	// 初始化日志
	InitLogger()
	// 初始化mysql
	InitMysql()
	// 初始化redis
	InitRedis()
	// 启动服务
	runServer()
}
