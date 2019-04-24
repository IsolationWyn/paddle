# 容器：提供一个与宿主机系统共享内核但与系统中的其它进程资源相隔离的执行环境

# builder
![](https://ws1.sinaimg.cn/large/006jGdC3gy1g2a1kzmpf5j313o18c483.jpg)

# Proc
![](https://ws1.sinaimg.cn/large/006jGdC3ly1g2256tpp5fj30su0f8jyg.jpg)

# Namespace
* UTS Namespace: 主机名隔离
* IPC Namespace: 进程间通信的隔离
* PID Namespace: 进程隔离
* Mount Namespace: 磁盘挂载点和文件系统的隔离
* User Namespace: 用户隔离
* Network Namespace: 网络栈隔离

# Linux 支持的subsystem
* cpu (since Linux 2.6.24; CONFIG_CGROUP_SCHED)     
用来限制cgroup的CPU使用率。
* cpuacct (since Linux 2.6.24; CONFIG_CGROUP_CPUACCT)           
统计cgroup的CPU的使用率。 
* cpuset (since Linux 2.6.24; CONFIG_CPUSETS)           
绑定cgroup到指定CPUs和NUMA节点。
* memory (since Linux 2.6.25; CONFIG_MEMCG)         
统计和限制cgroup的内存的使用率，包括process memory, kernel memory, 和swap。
* devices (since Linux 2.6.26; CONFIG_CGROUP_DEVICE)            
限制cgroup创建(mknod)和访问设备的权限。
* freezer (since Linux 2.6.28; CONFIG_CGROUP_FREEZER)               
suspend和restore一个cgroup中的所有进程。
* net_cls (since Linux 2.6.29; CONFIG_CGROUP_NET_CLASSID)           
将一个cgroup中进程创建的所有网络包加上一个classid标记，用于tc和iptables。 只对发出去的网络包生效，对收到的网络包不起作用。
* blkio (since Linux 2.6.33; CONFIG_BLK_CGROUP)         
限制cgroup访问块设备的IO速度。
* perf_event (since Linux 2.6.39; CONFIG_CGROUP_PERF)               
对cgroup进行性能监控
* net_prio (since Linux 3.3; CONFIG_CGROUP_NET_PRIO)                  
针对每个网络接口设置cgroup的访问优先级。
* hugetlb (since Linux 3.5; CONFIG_CGROUP_HUGETLB)              
限制cgroup的huge pages的使用量。
* pids (since Linux 4.3; CONFIG_CGROUP_PIDS)
限制一个cgroup及其子孙cgroup中的总进程数。

## Mount flags
* MS_BIND: 执行bind挂载, 使文件或者子目录在文件系统内的另一个点上可视
* MS_DIRSYNC: 同步目录的更新
* MS_MANDLOCK：允许在文件上执行强制锁。
* MS_MOVE：移动子目录树。
* MS_NOATIME：不要更新文件上的访问时间。
* MS_NODEV：不允许访问设备文件。
* MS_NODIRATIME：不允许更新目录上的访问时间。
* MS_NOEXEC：不允许在挂载的文件系统上执行程序。
* MS_NOSUID：执行程序时，不遵照set-user-ID 和 set-group-ID位。
* MS_RDONLY：指定文件系统为只读。
* MS_REMOUNT：重新加载文件系统。这允许你改变现存文件系统的mountflag和数据，而无需使用先卸载，再挂上文件系统方式。
* MS_SYNCHRONOUS：同步文件的更新。
* MNT_FORCE：强制卸载，即使文件系统处于忙状态。
* MNT_EXPIRE：将挂载点标志为过时。
* MS_STRICTATIME: 总是更新最后访问时间(始于linux 2.6.30)，定义了此,那么MS_NOATIME和MS_RELATIME将会忽略。
* MS_REC: 递归挂载, 跟MS_BIND配合


```go
type SysProcAttr struct {
	Chroot       string         // Chroot.
	Credential   *Credential    // Credential.
	Ptrace       bool           // Enable tracing.
	Setsid       bool           // Create session.
	Setpgid      bool           // Set process group ID to Pgid, or, if Pgid == 0, to new pid.
	Setctty      bool           // Set controlling terminal to fd Ctty (only meaningful if Setsid is set)
	Noctty       bool           // Detach fd 0 from controlling terminal
	Ctty         int            // Controlling TTY fd
	Foreground   bool           // Place child's process group in foreground. (Implies Setpgid. Uses Ctty as fd of controlling TTY)
	Pgid         int            // Child's process group ID if Setpgid.
	Pdeathsig    Signal         // Signal that the process will get when its parent dies (Linux only)
	Cloneflags   uintptr        // Flags for clone calls (Linux only)
	Unshareflags uintptr        // Flags for unshare calls (Linux only)
	UidMappings  []SysProcIDMap // User ID mappings for user namespaces.
	GidMappings  []SysProcIDMap // Group ID mappings for user namespaces.
	// GidMappingsEnableSetgroups enabling setgroups syscall.
	// If false, then setgroups syscall will be disabled for the child process.
	// This parameter is no-op if GidMappings == nil. Otherwise for unprivileged
	// users this should be set to false for mappings work.
	GidMappingsEnableSetgroups bool
	AmbientCaps                []uintptr // Ambient capabilities (Linux only)
}
```


* cgroup.clone_children

这个文件只对cpuset（subsystem）有影响，当该文件的内容为1时，新创建的cgroup将会继承父cgroup的配置，即从父cgroup里面拷贝配置文件来初始化新cgroup，可以参考这里

* cgroup.procs

当前cgroup中的所有进程ID，系统不保证ID是顺序排列的，且ID有可能重复

* cgroup.sane_behavior

* notify_on_release

该文件的内容为1时，当cgroup退出时（不再包含任何进程和子cgroup），将调用release_agent里面配置的命令。新cgroup被创建时将默认继承父cgroup的这项配置。

* release_agent

里面包含了cgroup退出时将会执行的命令，系统调用该命令时会将相应cgroup的相对路径当作参数传进去。 注意：这个文件只会存在于root cgroup下面，其他cgroup里面不会有这个文件。

* tasks

当前cgroup中的所有线程ID，系统不保证ID是顺序排列的

