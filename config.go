package ginapi

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func initConfig(configPath string) {
	fmt.Println("config loaded", configPath)
	if os.Getenv("mode") == "" {
		viper.Set("mode", "debug")
	}
	mode := viper.GetString("mode")
	viper.SetConfigFile(configPath + mode + ".yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic("配置文件读取失败，请检查 " + configPath + mode + ".yaml 是否存在")
	}
}
