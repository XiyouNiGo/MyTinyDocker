package main

import (
	"os"
	"strings"

	"github.com/XiyouNiGo/MyTinyDocker/cgroups"
	"github.com/XiyouNiGo/MyTinyDocker/cgroups/subsystems"
	"github.com/XiyouNiGo/MyTinyDocker/container"
	log "github.com/sirupsen/logrus"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	cgroupManager := cgroups.NewCgroupManager("mytinydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	parent.Wait()
	mntURL := "/root/mnt"
	rootURL := "/root"
	container.DeleteWorkSpace(rootURL, mntURL, volume)
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
