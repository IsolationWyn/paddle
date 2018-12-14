package container

import (
	log "github.com/Sirupsen/logrus"
	"syscall"
	"os/exec"
	"os"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	/*
	这里是父进程,也就是当前进程执行的内容
	1. 这里的/proc/self/exe 调用中, /proc/self指的是当前运行进程自己的环境, exec 其实就是调用了自己
	使用这种方式对创建出来的对象进程进行初始化
	2. 后面的args是参数, 其中init是传递给本进程的第一个参数, 在本例中, 其实就是会去调用initCommand去初始化进程的
	一些环境和资源
	3. 下面的clone参数就是去fork出来一个新进程, 并且使用了namespace隔离创建的进程和外部环境
	4. 如果用户指定了 -ti 参数, 就需要把当前进程的输入输出导入到标准输入输出上
	*/
	
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | 
				syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// 传入管道文件读取端的句柄
	// 一个进程默认有三个文件描述符(标准输入标准输出标准错误)
	cmd.ExtraFiles = []*os.File{readPipe}
	// cmd.Dir = "/root/busybox"
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	// 生成匿名管道
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil ,err
	}
	return read, write, nil
}

// NewWorkSpace函数是用来创建容器文件系统的, 它包括CreateReadOnlyLayer, CreateWriteLayer和CreateMountPoint
// CreateReadOnlyLayer函数新建busybox文件夹, 将busybox.tar解压到busybox目录下, 作为容器的只读层
// CreateWriteLayer函数创建一个名为writeLayer的文件夹, 作为容器唯一的可写层
// 在CreateMountPoint函数中, 首先创建了mnt文件夹, 作为挂载点, 然后把writeLayer目录和busybox目录mount到mnt目录下

func NewWorkSpace(rootURL string, mntURL string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
}

func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := rootURL + "busybox/"
	busyboxTarURL := rootURL + "busybox.tar"
	exist, err := PathExists(busyboxURL)
	if err != nil {
		log. Infof ("Fail to judge whether dir %s exists . %v", busyboxURL, err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			log.Errorf ("Mkdir dir %s error. %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", busyboxURL, err)
		}
	}
}

func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeURL, err)
	}
}

func CreateMountPoint(rootURL string, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntURL, err)
	}
	dirs := "dirs=" + rootURL + "writeLayer:" + rootURL + "busybox"
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}