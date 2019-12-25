package main

import (
	"crontab/worker"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	confFile string // 配置文件路径
)

//解析命令行参数
func initArgs() {
	//启动命令期望
	//worker -config ./worker.json
	//查看命令提示内容
	//worker -h
	//							参数名			默认值				提示
	flag.StringVar(&confFile, "config", "/Users/haoyong/workspace/src/crontab/worker/main/worker.json", "制定worker.json")
	//解析命令
	flag.Parse()
}

//初始化线程数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	var (
		err error
	)
	//初始化命令行参数
	initArgs()
	//初始化线程
	initEnv()
	//加载配置
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}
	//启动日志协程
	if err = worker.InitLogSink(); err != nil {
		goto ERR
	}
	//启动执行器
	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}
	//启动调度器
	if err = worker.InitScheduler(); err != nil {
		goto ERR
	}
	//任务管理器
	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	//正常退出
	for {
		time.Sleep(time.Second)
	}
	return
ERR:
	fmt.Println(err)
}
