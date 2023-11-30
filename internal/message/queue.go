package msgq

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type msgbuf struct {
	mtype int
	mtext []byte
}

// OpenQueue creates/opens a system v message queue
// and returns queue key and error, if any.
func OpenQueue(key int) (msqid uintptr, err error) {
	msqid, _, errno := syscall.Syscall(syscall.SYS_MSGGET, uintptr(key),
		unix.IPC_CREAT|unix.IPC_EXCL|0777, 0)
	if errno != 0 {
		err = fmt.Errorf("cannot create or open the message queue: %w", errno)
	}
	return
}

// DestroyQueue closes the queue has been opened
// returns error, if any.
func DestroyQueue(msqid uintptr) (err error) {
	_, _, errno := syscall.Syscall(syscall.SYS_MSGCTL, msqid, unix.IPC_RMID, 0)
	if errno != 0 {
		err = fmt.Errorf("cannot remove queue: %w", errno)
	}
	return
}

// SndMsg sends the messages to queue
// returns error if any.
func SndMsg(msqid uintptr, msgbf *msgbuf) (err error) {
	_, _, errno := syscall.Syscall(syscall.SYS_MSGSND, msqid,
		uintptr(unsafe.Pointer(msgbf)), uintptr(len(msgbf.mtext)))
	if errno != 0 {
		err = fmt.Errorf("cannot send message to queue: %w", errno)
	}
	return
}

// RcvMsg blocks and receive one message from queue.
func RcvMsg(msqid uintptr, msgbf *msgbuf) (err error) {
	for {
		_, _, errno := syscall.Syscall6(syscall.SYS_MSGRCV, msqid,
			uintptr(unsafe.Pointer(msgbf)), uintptr(len(msgbf.mtext)), uintptr(msgbf.mtype), 0, 0)
		if errno == syscall.EINTR {
			continue
		} else if errno != 0 {
      err = errno
		}
		return
	}
}

// NewMsg new a message.
func NewMsg(mtype int, mtext []byte) *msgbuf {
  if mtext == nil {
    mtext = []byte{}
  }
	return &msgbuf{mtype: mtype, mtext: mtext}
}

// MsgChan calls RcvMsg and returns a channel.
func MsgChan(msqid uintptr, msgbf *msgbuf) <-chan struct{} {
	ret := make(chan struct{})
	go func() {
		err := RcvMsg(msqid, msgbf)
		if err != nil {
      return
		}
		ret <- struct{}{}
	}()
	return ret
}
