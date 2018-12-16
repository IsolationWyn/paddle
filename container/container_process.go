package container

import (
	"strings"
	log "github.com/Sirupsen/logrus"
	"syscall"
	"os/exec"
	"os"
)

func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
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

	// 操作系统特定的创建属性, 用于控制进程中相关属性
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
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	NewWorkSpace(rootURL, mntURL, volume)
	cmd.Dir = mntURL
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

func NewWorkSpace(rootURL string, mntURL string, volume string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
	if (volume != "") {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if(length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "") {
			MountVolume(rootURL, mntURL, volumeURLs)
			log.Infof("%q", volumeURLs)
		} else {
			log.Infof("Volume parameter input is not correct.")
		}
	}
}

func volumeUrlExtract(volume string)([]string) {
	var volumeURLs []string
	volumeURLs = strings.Split(volume, ":")
	return volumeURLs
}

func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := rootURL + "/busybox"
	busyboxTarURL := rootURL + "/busybox.tar"
	exist, err := PathExists(busyboxURL)

	if err != nil {
		log. Infof ("Fail to judge whether dir %s exists . %v", busyboxURL, err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			log.Errorf ("Mkdir busybox dir %s error. %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", busyboxURL, err)
		}
	}
}

func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "/writeLayer"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Infof("Mkdir dir %s error. %v", writeURL, err)
	}
}

func CreateMountPoint(rootURL string, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		log.Infof("Mkdir mountpoint dir %s error. %v", mntURL, err)
	}
	dirs := "dirs=" + rootURL + "/writeLayer:" + rootURL + "/busybox"
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount mountpoint dir failed. %v", err)
	}
}

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		if os.IsNotExist(err) {
			return false, err
		}
	}
	return true, nil
}

// 挂载数据卷的过程
// 1. 首先, 读取宿主机文件目录URL, 创建宿主机文件目录(/root/${parentURL})
// 2. 然后, 读取容器挂载点URL, 在容器文件系统里创建挂载点(/root/mnt/${containerUrl})
// 3. 最后, 把宿主机文件目录挂载到容器挂载点, 这样启动容器

func MountVolume(rootURL string, mntURL string, volumeURLs []string) {
	// 创建宿主机文件目录
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil {
		log.Infof("Mkdir parent dir %s error. %v", parentUrl, err)
	}
	// 在容器文件系统里创建挂载点
	containerUrl := volumeURLs[1]
	containerVolumeURL := mntURL + containerUrl
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
		log.Infof("Mkdir container dir %s error. %v", containerVolumeURL, err)
	}
	// 把宿主机文件目录挂载到容器挂载点
	dirs := "dirs=" + parentUrl
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount volume failed. %v", err)
	}
}

// Docker 会在删除容器的时候, 把容器对应的Write Layer 和 Container-init Layer删除, 而保留镜像所有的内容。
// 首先, 在DeleteMountPoint函数中 umount mnt函数
// 然后, 删除mnt目录
// 最后, 在DeleteWriteLayer函数中删除writeLayer文件夹, 这样容器对文件系统的更改就都已经抹去了

func DeleteWorkSpace(rootURL string, mntURL string, volume string) {
	if(volume != "") {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if(length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "") {
			DeleteMountPointWithVolume(rootURL, mntURL, volumeURLs)
		} else {
			DeleteMountPoint(rootURL, mntURL)
		}
	}
	DeleteWriteLayer(rootURL)
}

func DeleteMountPointWithVolume(rootURL string, mntURL string, volumeURLs []string) {
	containerUrl := mntURL + volumeURLs[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Umount volume failed. %v", err)
	}
	
	cmd = exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Umount mountpoint failed. %v", err)
	}

	if err := os.RemoveAll(mntURL); err != nil {
		log.Infof("Remove mountpoint dir %s error %v", mntURL, err)
	}
}

func DeleteMountPoint(rootURL string, mntURL string) {
	cmd := exec.Command("unmount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("Remove dir %s error %v", mntURL, err)
	}
}

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "/writeLayer"
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("Remove writeLayer dir %s error %v", writeURL, err)
	}
}