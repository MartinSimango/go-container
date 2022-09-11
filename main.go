package main

import (
	"fmt"
	"os"
	"os/exec"

	// "os/signal"
	"io/ioutil"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("what??")
	}
}

func child() {
	fmt.Printf("running %v as %d\n", os.Args[2:], os.Getpid())
	ioutil.WriteFile("/sys/fs/cgroup/podrun/container/cgroup.procs", []byte(strconv.Itoa(os.Getpid())), 0777)
	ioutil.WriteFile("/sys/fs/cgroup/podrun/container/cpu.max", []byte("5000 100000"), 0777)
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot(("rootfs")))
	must(syscall.Chdir(("/")))
	// mount --make-private /
	// must(syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")) // do not propagte mnts to other mount namespaces
	must(syscall.Mount("none", "/proc", "proc", 0, ""))
	must(syscall.Mount("none", "/run/netns", "tmpfs", 0, ""))
	syscall.Clearenv()
	syscall.Setenv("HOME", "/")
	// syscall.Setenv("HOSTNAME", "container")

	// must(syscall.Mount("none", "/run/netns", "tmpfs", 0, ""))

	err := cmd.Run()
	exit_code := cmd.ProcessState.ExitCode()
	if exit_code == -1 {
		fmt.Println(err)
	}

}

func run() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	os.MkdirAll("/sys/fs/cgroup/podrun/container", 0700)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// signal.Ignore()
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	must(cmd.Run())
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
