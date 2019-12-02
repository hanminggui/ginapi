package common

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// FatalError 导致系统不能继续提供服务的错误
func FatalError(err error) {
	if err != nil {
		Mysql.Close()
		Redis.close()
		Log.Fatal(err, Stack(0))
	}
}

// PanicError panic error 让中间件去处理
func PanicError(err error) error {
	if err != nil {
		panic(err)
	}
	return err
}

// LogError 记录发生的错误信息
func LogError(c *gin.Context, err error) {
	if err != nil {
		logger, _ := c.Get("logger")
		entry := logger.(*logrus.Entry)
		entry.WithFields(logrus.Fields{
			"error": err,
			"stack": string(Stack(2)),
		}).Error("error")
	}
}

// CallBackError 如果出错运行 fn
func CallBackError(err error, fn func()) {
	if err != nil {
		fn()
	}
}
