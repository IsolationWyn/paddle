package main

import (
	"strings"
	"github.com/IsolationWyn/paddle/cgroups"
	"github.com/IsolationWyn/paddle/cgroups/subsystems"
	"github.com/IsolationWyn/paddle/container"
	log "github.com/Sirupsen/logrus"
	"os"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	/*
	这里的Start方法是真正开始前面创建好的command的调用, 它首先会
	clone出来一个namespace隔离进程, 然后在子进程中, 调用/proc/self/exe,
	也就是调用自己, 发送init参数, 调用init方法去初始化容器的一些资源
	*/
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	// 设置资源限制
	cgroupManager.Set(res)
	// 将容器进程加入到各个subsystem挂载对应的cgroup中
	cgroupManager.Apply(parent.Process.Pid)
	// 对容器设置完限制之后, 初始化容器
	sendInitCommand(comArray, writePipe)
	parent.Wait()
	os.Exit(-1)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}