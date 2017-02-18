package netns

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

type NsHandle int

// close closes the file descriptor mapped to a network namespace
func (h NsHandle) Close() error {
	return unix.Close(int(h))
}

func OpenNs(nsName string) (NsHandle, error) {
	fd, err := unix.Open(nsName, unix.O_RDONLY, 0)
	return NsHandle(fd), err
}

// setNs sets the process's network namespace
func SetNs(h NsHandle) error {

	return unix.Setns(int(h), unix.CLONE_NEWNET)
	//return unix.Setns(h, unix.CLONE_NEWNET)
}

func GetPath(path string) string {
	return fmt.Sprintf("/var/run/netns/%s", path)
}

func GetFromPath(path string) (NsHandle, error) {

	fd, err := unix.Open(fmt.Sprintf("/var/run/netns/%s", path),
		unix.O_RDONLY, 0)
	if err != nil {
		return -1, err
	}
	return NsHandle(fd), nil
}

func GetFromThread() (NsHandle, error) {
	return GetFromPath(fmt.Sprintf("/proc/%d/task/%d/ns/net", os.Getpid,
		unix.Gettid))
}
