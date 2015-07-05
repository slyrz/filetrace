package main

import "syscall"

func getSyscall(r *syscall.PtraceRegs) int {
	return int(r.Orig_rax)
}

func getArgAddr(r *syscall.PtraceRegs, argnum int) uintptr {
	switch argnum {
	case 0:
		return uintptr(r.Rdi)
	case 1:
		return uintptr(r.Rsi)
	case 2:
		return uintptr(r.Rdx)
	}
	return 0
}

func getArgInt(r *syscall.PtraceRegs, argnum int) int {
	switch argnum {
	case 0:
		return int(int32(r.Rdi & 0xffffffff))
	case 1:
		return int(int32(r.Rsi & 0xffffffff))
	case 2:
		return int(int32(r.Rdx & 0xffffffff))
	}
	return 0
}

var syscallTypes = map[int]int{
	syscall.SYS_ACCESS:     SyscallPath,
	syscall.SYS_CHDIR:      SyscallPath,
	syscall.SYS_CREAT:      SyscallPath,
	syscall.SYS_EXECVE:     SyscallPath,
	syscall.SYS_LCHOWN:     SyscallPath,
	syscall.SYS_LINK:       SyscallPath,
	syscall.SYS_LSTAT:      SyscallPath,
	syscall.SYS_MKDIR:      SyscallPath,
	syscall.SYS_OPEN:       SyscallPath,
	syscall.SYS_READLINK:   SyscallPath,
	syscall.SYS_RMDIR:      SyscallPath,
	syscall.SYS_STAT:       SyscallPath,
	syscall.SYS_STATFS:     SyscallPath,
	syscall.SYS_SYMLINK:    SyscallPath,
	syscall.SYS_TRUNCATE:   SyscallPath,
	syscall.SYS_UNLINK:     SyscallPath,
	syscall.SYS_UTIMES:     SyscallPath,
	syscall.SYS_FCNTL:      SyscallFile,
	syscall.SYS_FCHDIR:     SyscallFile,
	syscall.SYS_FACCESSAT:  SyscallFilePath,
	syscall.SYS_FCHMODAT:   SyscallFilePath,
	syscall.SYS_FCHOWNAT:   SyscallFilePath,
	syscall.SYS_LINKAT:     SyscallFilePath,
	syscall.SYS_MKDIRAT:    SyscallFilePath,
	syscall.SYS_MKNODAT:    SyscallFilePath,
	syscall.SYS_OPENAT:     SyscallFilePath,
	syscall.SYS_READLINKAT: SyscallFilePath,
	syscall.SYS_UNLINKAT:   SyscallFilePath,
}
