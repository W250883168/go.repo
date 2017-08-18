package xapp

import (
	"syscall"
	"unsafe"
)

var (
	NameMutex = "Shamem"
	kernel    = syscall.NewLazyDLL("kernel32.dll")
)

const (
	IPC_RMID   = 0
	IPC_CREAT  = 00001000
	IPC_EXCL   = 00002000
	IPC_NOWAIT = 00004000
)

func Lock_windows(Name string) (ptr uintptr, ok bool) {
	ptr, _, err := kernel.NewProc("CreateMutexA").Call(0, 1, uintptr(unsafe.Pointer(&NameMutex)))
	if err.Error() != "The operation completed successfully" {
		return ptr, ok
	}

	ok = true
	return ptr, ok
}

func UnLock_windows(id uintptr) {
	syscall.CloseHandle(syscall.Handle(id))
}

// 下面是linux系统上面的实现用的是使用共享内存的方法
//func Lock_linux(Name string) (uintptr, bool) {
//	id, _, err := syscall.Syscall6(syscall.SYS_SHMGET, uintptr(unsafe.Pointer(&NameMutex)), 1, IPC_CREAT|IPC_EXCL, 0, 0, 0)
//	if err.Error() != "errno 0" {
//		return 0, false
//	}
//	return id, true
//}

//func UnLock_linux(id uintptr) {
//	syscall.Syscall6(syscall.SYS_SHMCTL, id, IPC_RMID, 0, 0, 0, 0)
//}
