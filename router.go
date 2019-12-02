package ginapi

import (
	"fmt"
	. "github.com/hanminggui/ginapi/log"
	. "github.com/hanminggui/ginapi/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"reflect"
	"strings"
)

var (
	funcs map[string]map[string]*reflect.Value
	// apis 所有需要映射的api结构体都放进来
)

/**
 * 初始化
 * 反射apis下面的所有func 包名/funcName 存到FUNCS map
 * 加载api到路由
 */
func initRouter(engine *gin.Engine) {

	Log.Debug("init reflect routers.")
	Log.Debug("apiStucts count:", len(apiStucts))
	funcs = make(map[string]map[string]*reflect.Value)
	for _, m := range []string{"GET", "POST", "PUT", "PATCH", "HEAD", "OPTIONS", "DELETE"} {
		funcs[m] = make(map[string]*reflect.Value)
	}
	for _, element := range apiStucts {
		eValue := reflect.ValueOf(element)
		eType := eValue.Type()
		for i := 0; i < eValue.NumMethod(); i++ {
			methodName := eType.Method(i).Name
			oldHead := string(methodName)
			newHead := strings.ToLower(oldHead)
			methodName = strings.Replace(methodName, oldHead, newHead, 1)
			// Log.Info(strings.SplitN(eType.PkgPath(), viper.GetString("product"), 2))
			apiPath := "/" + strings.SplitN(eType.PkgPath(), viper.GetString("product"), 2)[0] + "/" + methodName
			fn := eValue.Method(i)
			method := eType.Name()
			funcs[method][apiPath] = &fn
			Log.Debugf("load api	%s	-->	%s", method, apiPath)
		}
	}

	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	lodaApis(engine)
}

/**
 * 请求api后包装返回结果
 */
func call(c *gin.Context, fn *reflect.Value) {
	paramList := []reflect.Value{reflect.ValueOf(c)}
	Log.Info(3, fn)
	retList := fn.Call(paramList)
	switch resType := retList[0].Interface().(type) {
	case error:
		UnknownError(c, resType.Error())
	default:
		Success(c, retList[0].Interface())
	}
}

/**
 * 加载api到路由
 */
func lodaApis(engine *gin.Engine) {
	for path := range funcs["GET"] {
		engine.GET(path, func(c *gin.Context) {
			Log.Info(0, funcs["GET"])
			Log.Info(1, fmt.Sprint(c.Request.URL.Path))
			Log.Info(2, funcs["GET"][fmt.Sprint(c.Request.URL.Path)])
			call(c, funcs["GET"][fmt.Sprint(c.Request.URL.Path)])
		})
	}
	for path := range funcs["POST"] {
		engine.POST(path, func(c *gin.Context) {
			call(c, funcs["POST"][fmt.Sprint(c.Request.URL.Path)])
		})
	}
	for path := range funcs["PUT"] {
		engine.PUT(path, func(c *gin.Context) {
			call(c, funcs["PUT"][fmt.Sprint(c.Request.URL.Path)])
		})
	}
	for path := range funcs["PATCH"] {
		engine.PATCH(path, func(c *gin.Context) {
			call(c, funcs["PATCH"][fmt.Sprint(c.Request.URL.Path)])
		})
	}
	for path := range funcs["HEAD"] {
		engine.HEAD(path, func(c *gin.Context) {
			call(c, funcs["HEAD"][fmt.Sprint(c.Request.URL.Path)])
		})
	}
	for path := range funcs["OPTIONS"] {
		engine.OPTIONS(path, func(c *gin.Context) {
			call(c, funcs["OPTIONS"][fmt.Sprint(c.Request.URL.Path)])
		})
	}
	for path := range funcs["DELETE"] {
		engine.DELETE(path, func(c *gin.Context) {
			call(c, funcs["DELETE"][fmt.Sprint(c.Request.URL.Path)])
		})
	}
}
