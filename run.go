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
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 创建cgroup manager, 并通过调用set和apply设置资源限制并使限制在容器上生效
	cgroupManager := cgroups.NewCgroupManager("paddle")
	defer cgroupManager.Destroy()
	// 设置资源限制
	cgroupManager.Set(res)
	// 将容器进程加入到各个subsystem挂载对应的cgroup中
	cgroupManager.Apply(parent.Process.Pid)
	// 对容器设置完限制之后, 初始化容器
	sendInitCommand(comArray, writePipe)
	parent.Wait()
	// os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}