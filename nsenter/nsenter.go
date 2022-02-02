package nsenter

/*
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>

__attribute__((constructor)) void enter_namespace(void) {
	char *mytinydocker_pid;
	mytinydocker_pid = getenv("mytinydocker_pid");
	if (mytinydocker_pid) {
		//fprintf(stdout, "got mytinydocker_pid=%s\n", mytinydocker_pid);
	} else {
		//fprintf(stdout, "missing mytinydocker_pid env skip nsenter");
		return;
	}
	char *mytinydocker_cmd;
	mytinydocker_cmd = getenv("mytinydocker_cmd");
	if (mytinydocker_cmd) {
		//fprintf(stdout, "got mytinydocker_cmd=%s\n", mytinydocker_cmd);
	} else {
		//fprintf(stdout, "missing mytinydocker_cmd env skip nsenter");
		return;
	}
	int i;
	char nspath[1024];
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };

	for (i=0; i<5; i++) {
		sprintf(nspath, "/proc/%s/ns/%s", mytinydocker_pid, namespaces[i]);
		int fd = open(nspath, O_RDONLY);

		if (setns(fd, 0) == -1) {
			//fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		} else {
			//fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
		}
		close(fd);
	}
	int res = system(mytinydocker_cmd);
	exit(0);
	return;
}
*/
import "C"
