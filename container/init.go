package container

import (
	"path/filepath"
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
	
	// defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	setupMount()

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


func pivotRoot(root string) error {
	/**
	为了使当前root和老root不在同一文件系统下, 重新mount root
	bind mount 是把相同的内容换了一个挂载点的方法
	*/
	// 执行bind挂载, 使文件或者子目录在文件系统内的另一个点上可视 | 
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}

	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %var", err)
	}

	// 删除临时文件夹
	return os.Remove(pivotDir)
}


/*
init 挂载点
*/
func setupMount() {
	// 获取当前路径
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return 
	}
	log.Infof("Current location is %s", pwd)
	pivotRoot(pwd)

	// mount proc
	// 不允许在挂载的文件系统上执行程序 | 执行程序时，不遵照set-user-ID 和 set-group-ID位 | 不允许访问设备文件
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	// 执行程序时，不遵照set-user-ID 和 set-group-ID位 | 总是更新最后访问时间
	// tmpfs 是一种基于内存的文件系统, 可以使用RAM或swap分区来存储
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}