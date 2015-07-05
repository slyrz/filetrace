package main

const (
	// ignore syscall
	SyscallIgnore = iota
	// function of type f(char *path, ...)
	SyscallPath
	// function of type f(int fd, ...)
	SyscallFile
	// function of type f(int fd, char *path, ...)
	SyscallFilePath
)

