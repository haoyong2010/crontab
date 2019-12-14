package main

import (
	"crontab/master"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	confFile string //配置文件路径
)

//解析命令行参数
func initArgs() {
	//启动命令期望
	//master -config ./master.json
	//查看命令提示内容
	//master -h
	//						 参数名		     默认值				   提示
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
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
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	//任务管理器
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	//启动API HTTP服务
	if err = master.InitApiServer(); err != nil {
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
