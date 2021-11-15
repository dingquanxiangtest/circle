package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/molecule/api/restful"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
)

var (
	configPath = flag.String("config", "../configs/config.yml", "-config 配置文件地址")
)

func main() {
	flag.Parse()

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}

	err = logger.New(&conf.Log)
	if err != nil {
		panic(err)
	}

	// 启动路由
	router, err := restful.NewRouter(conf)
	if err != nil {
		panic(err)
	}
	go router.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			router.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
