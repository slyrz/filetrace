package main

import "syscall"

var syscallType = map[int]int{
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

var syscallCwdChange = map[int]bool{
	syscall.SYS_CHDIR:  true,
	syscall.SYS_FCHDIR: true,
}
