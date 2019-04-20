package nsenter

/*
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
	paddle_pid = getenv("paddle_pid");
	if (paddle_pid) {
		//fprintf(stdout, "got paddle_pid=%s\n", paddle_pid);
	} else {
		//fprintf(stdout, "missing paddle_pid env skip nsenter");
		return;
	}
	char *paddle_cmd;
	paddle_cmd = getenv("paddle_cmd");
	if (paddle_cmd) {
		//fprintf(stdout, "got paddle_cmd=%s\n", paddle_cmd);
	} else {
		//fprintf(stdout, "missing paddle_cmd env skip nsenter");
		return;
	}
	int i;
	char nspath[1024];
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };
	for (i=0; i<5; i++) {
		sprintf(nspath, "/proc/%s/ns/%s", paddle_pid, namespaces[i]);
		int fd = open(nspath, O_RDONLY);
		if (setns(fd, 0) == -1) {
			//fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		} else {
			//fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
		}
		close(fd);
	}
	int res = system(paddle_cmd);
	exit(0);
	return;
}
*/
import "C"