package sleep

import "syscall"

func NanoSleep(usec int64) error {
	req := syscall.Timespec{Sec: 0, Nsec: 1e6 * usec}
	for {
		if err := syscall.Nanosleep(&req, &req); err == syscall.EINTR {
			continue
		} else {
			return err
		}
	}
}
