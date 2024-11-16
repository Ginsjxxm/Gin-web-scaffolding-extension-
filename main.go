package main

import (
	"BlueBell/dao/mysql"
	"BlueBell/dao/redis"
	"BlueBell/logger"
	routers "BlueBell/routes"
	"BlueBell/settings"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//加载配置
	if err := settings.Init(); err != nil {
		fmt.Println("init setting failed,err:", err)
		return
	}
	//加载日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Println("init logger failed,err:", err)
		return
	}
	zap.L().Info("logger init success")
	defer zap.L().Sync()

	//加载数据库
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Println("init setting failed,err:", err)
		return
	}
	defer mysql.Close()

	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("init redis failed,err:", err)
		return
	}
	defer redis.Close()

	r := routers.SetupRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutdown Server ...")

	// 创建上下文，设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置超时时间为 10 秒
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown failed", zap.Error(err))
	}

	zap.L().Info("Server exiting")

}
