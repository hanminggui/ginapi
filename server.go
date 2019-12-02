package ginapi

import (
	"context"
	. "ginapi/log"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func runServer() {
	gin.SetMode(viper.GetString("mode"))
	r := gin.New()
	initMiddlewares(r)
	initRouter(r)
	srv := &http.Server{
		Addr:    ":" + viper.GetString("service.port"),
		Handler: r,
	}
	// go r.Run()
	go func() {
		// 服务连接
		Log.Infof("listen and server %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	Log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		Log.Fatal("Server Shutdown:", err)
	}
	Log.Println("Server exiting")
}
