# Practise with <<Write My Own Docker>>

## Namespace
* UTS Namespace: 主机名隔离
* IPC Namespace: 进程间通信的隔离
* PID Namespace: 进程隔离
* Mount Namespace: 磁盘挂载点和文件系统的隔离
* User Namespace: 用户隔离
* Network Namespace: 网络栈隔离

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