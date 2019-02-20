#define _GNU_SOURCE
#include <unistd.h>
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>


__attribute__((constructor)) void enter_namespace(void) {
	char *paddle_pid;
   // 从环境变量获取需要进入的pid
	paddle_pid = getenv("paddle_pid");
	if (paddle_pid) {
		//fprintf(stdout, "got mydocker_pid=%s\n", mydocker_pid);
	} else {
		//fprintf(stdout, "missing mydocker_pid env skip nsenter");
		return;
	}
   char *paddle_cmd;
   // 从环境变量里面执行需要执行的命令
   paddle_cmd = getenv("paddle_cmd");
   if (paddle_cmd) {
		//fprintf(stdout, "got mydocker_cmd=%s\n", mydocker_cmd);
	} else {
		//fprintf(stdout, "missing mydocker_cmd env skip nsenter");
		return;
	}
   int i;
   char nspath[1024];
   char *namespaces[] = {"ipc", "uts", "net", "pid", "mnt"};

   for (i=0; i<5; i++) {
      // 拼接类似于/proc/pid/ns/ipc
      sprintf(nspath, "/proc/%s/ns/%s", paddle_pid, namespaces[i]);
      int fd = open(nspath, O_RDONLY);
      // 调用setns系统调用进入对应的Namespace
      if (setns(fd, 0) == -1) {
         //fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
      } else {
         //fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
      }
      close(fd);
   }
   // 在进入的Namespace中执行指定的命令
   int res = system(paddle_cmd);
   exit(0);
   return;
}