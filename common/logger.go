package common

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

// Log 日志实例
var Log *logrus.Logger

func InitLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}
	Log = logrus.New()
	if viper.GetString("mode") != "relese" {
		Log.SetLevel(logrus.DebugLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: viper.GetString("log.timeFormat"),
			FullTimestamp:   true,
		})
	} else { // 生产环境log模式
		Log.SetLevel(logrus.InfoLevel)
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Log.Info("关闭控制台log失败")
			Log.SetOutput(src)
		}
		apiLogPath := viper.GetString("log.fileName")
		logWriter, err := rotatelogs.New(
			apiLogPath+viper.GetString("log.splitName"),
			rotatelogs.WithLinkName(apiLogPath),                                                   // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(time.Duration(viper.GetInt64("log.maxHour"))*time.Hour),         // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Duration(viper.GetInt64("log.splitHour"))*time.Hour), // 日志切割时间间隔
		)
		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
		}
		lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
		Log.AddHook(lfHook)
		Log.SetLevel(logrus.InfoLevel)
	}
	return Log
}
