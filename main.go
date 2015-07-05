package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"sort"
	"strings"
)

const (
	EventFork  = 1
	EventVFork = 2
	EventClone = 3
	EventExec  = 4
)

var ignore = []string{
	"/tmp",
	"/proc",
}

func waitEvent(status syscall.WaitStatus) uint32 {
	return (uint32(status)>>16) & 0xff
}

func wait(pid int) (cpid int, status syscall.WaitStatus, err error) {
	cpid, err = syscall.Wait4(pid, &status, syscall.WALL, nil)
	if err != nil {
		return 0, 0, err
	}
	return cpid, status, nil
}

// traceSyscalls executes a command exec with the arguments args and calls
// the function callback for every syscall executed by the command's process.
func traceSyscalls(proc *os.Process, state func(pid int, regs *syscall.PtraceRegs)) {
	// flags used with ptrace
	const flags = syscall.PTRACE_O_TRACEVFORK |
		syscall.PTRACE_O_TRACEFORK |
		syscall.PTRACE_O_TRACECLONE |
		syscall.PTRACE_O_TRACEEXEC |
		syscall.PTRACE_O_TRACESYSGOOD

	if err := syscall.PtraceSetOptions(proc.Pid, flags); err != nil {
		log.Fatal("PtrageSetOptions", err)
	}
	if err := syscall.PtraceSyscall(proc.Pid, 0); err != nil {
		log.Fatalf("PtraceCont: %v", err)
	}
	for {
		signal := 0
		pid, status, err := wait(-1)
		if err != nil {
			break
		}
		if status.Exited() || status.Signaled() {
			continue
		}
		if status.Stopped() {
			switch waitEvent(status) {
			case EventFork, EventVFork, EventClone, EventExec:
				if cpid, err := syscall.PtraceGetEventMsg(pid); err == nil {
					log.Printf("process %d created new process %d", pid, cpid)
				}
			default:
				if stopSignal := status.StopSignal(); stopSignal&0x7f != syscall.SIGTRAP {
					signal = int(stopSignal)
				}
				var regs syscall.PtraceRegs
				if syscall.PtraceGetRegs(pid, &regs) == nil {
					state(pid, &regs)
				}
			}
		}
		syscall.PtraceSyscall(pid, signal)
	}
}

// getString returns the C string stored at the process' memory address addr.
func getString(pid int, addr uintptr) string {
	var buffer [4096]byte
	if _, err := syscall.PtracePeekData(pid, addr, buffer[:]); err == nil {
		if i := bytes.IndexByte(buffer[:], 0); i >= 0 && i < len(buffer) {
			return string(buffer[:i])
		}
	}
	return ""
}

// getLink returns the path to which the process' openend file descriptor
// fd belongs.
func getLink(pid int, fd int) string {
	if link, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%d", pid, fd)); err == nil {
		return link
	}
	return ""
}

// getCwd returns the process' current working directory.
func getCwd(pid int) string {
	if link, err := os.Readlink(fmt.Sprintf("/proc/%d/cwd", pid)); err == nil {
		return link
	}
	return ""
}

// inside returs true if path lies inside dir.
func inside(path, dir string) bool {
	path = path + "/"
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	return strings.HasPrefix(path, dir)
}

// ignorePath returs true if path lies in any directory of the ignore list.
func ignorePath(path string) bool {
	for _, dir := range ignore {
		if inside(path, dir) {
			return true
		}
	}
	return false
}

func main() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()

	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "Call: %s COMMAND\n", os.Args[0])
		return
	}

	proc, err := os.StartProcess(os.Args[1], os.Args[1:], &os.ProcAttr{
		Dir:   ".",
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Sys: &syscall.SysProcAttr{
			Ptrace:    true,
			Pdeathsig: syscall.SIGCHLD,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_, status, err := wait(proc.Pid)
	if err != nil {
		log.Fatalf("wait: %v", err)
	}
	if !status.Stopped() || status.StopSignal() != syscall.SIGTRAP {
		log.Fatalf("status: got %#x, want %#x", status, 0x57f)
	}

	paths := make(map[string]bool)

	// Step 1: add the process executable path.
	if link, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", proc.Pid)); err == nil {
		paths[link] = true
	}

	// Step 2: trace syscalls and add all used paths.
	traceSyscalls(proc, func(pid int, regs *syscall.PtraceRegs) {
		call, work, path := getSyscall(regs), "", ""
		switch syscallType[call] {
		case SyscallPath:
			path = getString(pid, getArgAddr(regs, 0))
		case SyscallFile:
			path = getLink(pid, getArgInt(regs, 0))
		case SyscallFilePath:
			work = getLink(pid, getArgInt(regs, 0))
			path = getString(pid, getArgAddr(regs, 1))
		default:
			return
		}
		if path == "" || path == "." || path == ".." {
			return
		}
		if !filepath.IsAbs(path) {
			if work == "" {
				work = getCwd(pid)
			}
			path = filepath.Join(work, path)
		}
		log.Printf("%q", path)
		paths[path] = true
	})

	// Remove paths that do not exist or shall be ignored.
	for path := range paths {
		if _, err := os.Stat(path); err != nil || ignorePath(path) {
			log.Printf("ignoring %q", path)
			delete(paths, path)
		}
	}

	// Sort paths to make the output look nice.
	list := make([]string, 0, len(paths))
	for path := range paths {
		list = append(list, path)
	}
	sort.StringSlice(list).Sort()

	for _, path := range list {
		fmt.Println(path)
	}
}
