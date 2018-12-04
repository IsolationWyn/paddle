package container

import (
	"os/exec"
	"strings"
	"io/ioutil"
	"fmt"
	"os"
	"syscall"
	log "github.com/Sirupsen/logrus"
)

func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}
	/*
	使用mount去挂载proc文件系统, 以便后面通过ps等命令去查看当前进程资源的情况
	init进程读取了父进程传递过来的参数后, 在子进程内进行了执行, 这样就完成了将用户指定命令传递给子进程的操作

	*/
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	// 调用exec.LookPath, 可以在系统的PATH里面寻找命令的绝对路径
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error", err)
		return err
	} 
	log.Infof("Find path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}


