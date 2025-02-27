package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/klog/v2"

	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/nodeinfodb"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/server"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/tasks"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/tasks/command"
	dockerinstall "github.com/edgewize-io/edgewize/edgewize-ota/pkg/tasks/docker-install"
)

func main() {
	dbPath := flag.String("db", "nodeinfo.db", "sqlite db path")
	innerDir := flag.String("innerDir", "./roles", "Directory for inner tasks")
	outerDir := flag.String("outerDir", "./tasks", "Directory for outer tasks")
	// 加载日志配置
	// 初始化 klog
	klog.InitFlags(flag.CommandLine)

	flag.Parse()
	nodeinfodb.InitDB(*dbPath)

	tasks.InitTasks(*innerDir, *outerDir)
	command.Init()
	dockerinstall.Init()

	// 提供 http 接口接口服务
	go server.Run()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	// 等待信号
	<-signalCh
}
