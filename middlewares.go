package ginapi

import (
	"bytes"
	"encoding/json"
	. "github.com/hanminggui/ginapi/common"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// InitMiddlewares 初始化中间件
func initMiddlewares(engine *gin.Engine) {
	engine.NoMethod(NotFound)
	engine.NoRoute(NotFound)
	engine.Use(requestID, logger, token2MemID, recovery)
}

// RequestID 生成请求id
func requestID(c *gin.Context) {
	requestID := uuid.NewV4().String()
	c.Set("requestId", requestID)
	c.Next()
}

// logger 记录请求日志
func logger(c *gin.Context) {
	entry := Log.WithField("requestId", c.GetString("requestId"))
	c.Set("logger", entry)
	data := "{}"
	if c.Request.Method == "POST" {
		tmpData, err := c.GetRawData()
		if err != nil {
			entry.Warning(err.Error())
		}
		data = string(tmpData)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(tmpData))
	}
	entry.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"ip":     c.ClientIP(),
		"url":    c.Request.URL,
		"proto":  c.Request.Proto,
		"data":   data,
	}).Info("request")

	beginTime := time.Now()
	c.Next()
	endTime := time.Now()

	res, e := c.Get("response")
	var b []byte
	if e {
		mp := res.(*APIResponse)
		b, _ = json.Marshal(mp)
	}
	entry.WithFields(logrus.Fields{
		"res_size": c.Writer.Size(),
		"status":   c.Writer.Status(),
		"use_time": endTime.Sub(beginTime),
		"response": string(b),
	}).Info("response")
}

func token2MemID(c *gin.Context) {
	token := c.GetHeader("token")
	var memID uint
	exists := Redis.Get(token, &memID)
	if exists {
		c.Set("memId", memID)
	}
	c.Next()
}

// recovery 捕获pance
func recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger, _ := c.Get("logger")
			entry := logger.(*logrus.Entry)
			entry.WithFields(logrus.Fields{
				"error": err,
				"stack": string(Stack(4)),
			}).Error("error")
			ServerError(c)
		}
	}()
	c.Next()
}
