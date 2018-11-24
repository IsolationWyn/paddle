package main

import (
	"github.com/IsolationWyn/paddle/container"
	log "github.com/Sirupsen/logrus"
	"os"
)

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	/*
	这里的Start方法是真正开始前面创建好的command的调用, 它首先会
	clone出来一个namespace隔离进程, 然后在子进程中, 调用/proc/self/exe,
	也就是调用自己, 发送init参数, 调用init方法去初始化容器的一些资源
	*/
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
