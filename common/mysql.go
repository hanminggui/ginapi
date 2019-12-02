package common

import (
	"fmt"
	. "github.com/hanminggui/ginapi/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"time"
)

// Mysql *gorm.DB实例
var Mysql *gorm.DB

func InitMysql() {
	if Mysql != nil {
		return
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "test_" + defaultTableName
	}
	config := viper.GetStringMapString("mysql")
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", config["user"], config["password"], config["addr"], config["database"])
	database, err := gorm.Open("mysql", dns)
	FatalError(err)
	Mysql = database
	Mysql.SetLogger(Log)
	Mysql.LogMode(true)
	Mysql.DB().SetMaxIdleConns(viper.GetInt("mysql.maxIdle"))
	Mysql.DB().SetMaxOpenConns(viper.GetInt("mysql.maxOpen"))
	Mysql.DB().Ping()
	Mysql.SingularTable(true)
	go heartbeatMysql()
}

// 心跳检测
func heartbeatMysql() {
	pingMysql()
	t := time.Tick(time.Second * 10)
	for range t {
		pingMysql()
	}
}

func pingMysql() {
	FatalError(Mysql.DB().Ping())
	Log.Debug("ping mysql success")
}
