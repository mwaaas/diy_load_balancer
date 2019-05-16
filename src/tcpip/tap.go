// +build linux

package tcpip

import (
	"errors"
	"log"
	"os/exec"
	"syscall"
	"unsafe"
)

var tap *Tap

type Tap struct {
	Name string
	Fd   int
}

func NewTap(name string) *Tap {
	tap = &Tap{Name: name}
	return tap
}

// Open opens the specified TUN device, and
// returns its file descriptor.
func (t *Tap) Open() error {
	fd, err := syscall.Open("/dev/net/tun", syscall.O_RDWR, 0)

	if err != nil {
		log.Fatalf("error opening tap device: %s", err)
	}

	t.Fd = fd

	var ifr struct {
		name  [16]byte
		flags uint16
		//_ [22]byte
	}

	copy(ifr.name[:], t.Name)

	/* Flags: IFF_TUN   - TUN device (no Ethernet headers)
	 *        IFF_TAP   - TAP device
	 *
	 *        IFF_NO_PI - Do not provide packet information
	 */
	ifr.flags = syscall.IFF_TAP | syscall.IFF_NO_PI

	_, _, erron := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(fd),
		syscall.TUNSETIFF,
		uintptr(unsafe.Pointer(&ifr)))

	if erron != 0 {
		syscall.Close(fd)
		return erron
	}

	t.Fd = fd
	return nil
}

func (t *Tap) Read(b []byte) (int, error) {
	return syscall.Read(t.Fd, b)
}

func (t *Tap) Write(b []byte) (int, error) {
	return syscall.Write(t.Fd, b)
}

func (t *Tap) Close() error {
	return syscall.Close(t.Fd)
}

func (t *Tap) SetAddress(addr string) error {
	info, err := exec.Command("ip", "addr", "add", addr, "dev", t.Name).CombinedOutput()
	if err != nil {
		return errors.New(err.Error() + " " + string(info))
	}
	return nil
}

//  ("ip link set dev %s up", dev);
func (t *Tap) SetUp() error {
	info, err := exec.Command("ip", "link", "set", "dev", t.Name, "up").CombinedOutput()
	if err != nil {
		return errors.New(err.Error() + " " + string(info))
	}
	return nil
}
