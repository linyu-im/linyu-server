package main

import (
	"flag"
	"fmt"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/email"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/i18n"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/logger"
	"github.com/linyu-im/linyu-server/linyu-gateway/pkg/gateway"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configPath = flag.String("config", "assets/config/config.yml", "Path to configuration file")
var localesDir = flag.String("locales", "assets/locales", "Directory for i18n locale files")
var emailDir = flag.String("email-templates", "assets/email_templates", "Directory for email templates")

func main() {
	flag.Parse()
	config.InitConfig(*configPath) //配置初始化
	logger.InitLog()               //日志初始化
	db.InitDB()                    //数据库相关初始化
	email.InitEmail(*emailDir)     //邮件相关初始化
	i18n.InitI18n(*localesDir)     //国际化初始化
	gateway.Run()                  //服务运行

	defer func() {
		//日志缓冲区内容强制刷新
		if err := logger.Log.Sync(); err != nil {
			fmt.Println("log sync error:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	// 阻塞等待信号
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logger.Log.Info("Shutting down server...")
			// 清理资源内容...
			time.Sleep(1 * time.Second)
			logger.Log.Info("Server stopped gracefully")
			return
		case syscall.SIGHUP:
			logger.Log.Info("Received SIGHUP, ignoring...")
		default:
			logger.Log.Info("Received other signal, exiting")
			return
		}
	}
}
