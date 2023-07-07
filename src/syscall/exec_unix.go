// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

// Fork, exec, wait, etc.

package syscall

import (
	errorspkg "errors"
	"internal/bytealg"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// ForkLock is used to synchronize creation of new file descriptors
// with fork.
//
// We want the child in a fork/exec sequence to inherit only the
// file descriptors we intend. To do that, we mark all file
// descriptors close-on-exec and then, in the child, explicitly
// unmark the ones we want the exec'ed program to keep.
// Unix doesn't make this easy: there is, in general, no way to
// allocate a new file descriptor close-on-exec. Instead you
// have to allocate the descriptor and then mark it close-on-exec.
// If a fork happens between those two events, the child's exec
// will inherit an unwanted file descriptor.
//
// This lock solves that race: the create new fd/mark close-on-exec
// operation is done holding ForkLock for reading, and the fork itself
// is done holding ForkLock for writing. At least, that's the idea.
// There are some complications.
//
// Some system calls that create new file descriptors can block
// for arbitrarily long times: open on a hung NFS server or named
// pipe, accept on a socket, and so on. We can't reasonably grab
// the lock across those operations.
//
// It is worse to inherit some file descriptors than others.
// If a non-malicious child accidentally inherits an open ordinary file,
// that's not a big deal. On the other hand, if a long-lived child
// accidentally inherits the write end of a pipe, then the reader
// of that pipe will not see EOF until that child exits, potentially
// causing the parent program to hang. This is a common problem
// in threaded C programs that use popen.
//
// Luckily, the file descriptors that are most important not to
// inherit are not the ones that can take an arbitrarily long time
// to create: pipe returns instantly, and the net package uses
// non-blocking I/O to accept on a listening socket.
// The rules for which file descriptor-creating operations use the
// ForkLock are as follows:
//
//   - Pipe. Use pipe2 if available. Otherwise, does not block,
//     so use ForkLock.
//   - Socket. Use SOCK_CLOEXEC if available. Otherwise, does not
//     block, so use ForkLock.
//   - Open. Use O_CLOEXEC if available. Otherwise, may block,
//     so live with the race.
//   - Dup. Use F_DUPFD_CLOEXEC or dup3 if available. Otherwise,
//     does not block, so use ForkLock.
var ForkLock sync.RWMutex

// StringSlicePtr converts a slice of strings to a slice of pointers
// to NUL-terminated byte arrays. If any string contains a NUL byte
// this function panics instead of returning an error.
//
// Deprecated: Use SlicePtrFromStrings instead.
func StringSlicePtr(ss []string) []*byte {
	bb := make([]*byte, len(ss)+1)
	for i := 0; i < len(ss); i++ {
		bb[i] = StringBytePtr(ss[i])
	}
	bb[len(ss)] = nil
	return bb
}

// SlicePtrFromStrings converts a slice of strings to a slice of
// pointers to NUL-terminated byte arrays. If any string contains
// a NUL byte, it returns (nil, EINVAL).
func SlicePtrFromStrings(ss []string) ([]*byte, error) {
	n := 0
	for _, s := range ss {
		if bytealg.IndexByteString(s, 0) != -1 {
			return nil, EINVAL
		}
		n += len(s) + 1 // +1 for NUL
	}
	bb := make([]*byte, len(ss)+1)
	b := make([]byte, n)
	n = 0
	for i, s := range ss {
		bb[i] = &b[n]
		copy(b[n:], s)
		n += len(s) + 1
	}
	return bb, nil
}

func CloseOnExec(fd int) { fcntl(fd, F_SETFD, FD_CLOEXEC) }

func SetNonblock(fd int, nonblocking bool) (err error) {
	flag, err := fcntl(fd, F_GETFL, 0)
	if err != nil {
		return err
	}
	if nonblocking {
		flag |= O_NONBLOCK
	} else {
		flag &^= O_NONBLOCK
	}
	_, err = fcntl(fd, F_SETFL, flag)
	return err
}

// Credential holds user and group identities to be assumed
// by a child process started by StartProcess.
type Credential struct {
	Uid         uint32   // User ID.
	Gid         uint32   // Group ID.
	Groups      []uint32 // Supplementary group IDs.
	NoSetGroups bool     // If true, don't set supplementary groups
}

// ProcAttr holds attributes that will be applied to a new process started
// by StartProcess.
type ProcAttr struct {
	Dir   string    // Current working directory.
	Env   []string  // Environment.
	Files []uintptr // File descriptors.
	Sys   *SysProcAttr
}

var zeroProcAttr ProcAttr
var zeroSysProcAttr SysProcAttr

type baseProcessHandle struct {
	terminated atomic.Bool
	holdMu     sync.RWMutex
}

func (ph *baseProcessHandle) markTerminated(graceful bool) {
	ph.terminated.Store(true)
	if graceful {
		// Acquire a write lock on holdMu to wait for any tasks
		// that hold the process (e.g., signalling) complete.
		ph.holdMu.Lock()
		ph.holdMu.Unlock()
	}
}

func (ph *baseProcessHandle) Hold() bool {
	ph.holdMu.RLock()
	if ph.terminated.Load() {
		ph.holdMu.RUnlock()
		return false
	}
	return true
}

func (ph *baseProcessHandle) Unhold() {
	ph.holdMu.RUnlock()
}

// managedProcessHandle is created for processes for which Wait4()
// should not be called directly. Instead, wait results will be
// delivered by a dedicated goroutine that calls Wait4(-1). This is
// necessary to ensure that all zombie processes are removed when
// running as PID 1.
type managedProcessHandle struct {
	baseProcessHandle

	haveStatus chan struct{}
	err        error
	status     WaitStatus
	rusage     Rusage
}

func (ph *managedProcessHandle) Wait4(status *WaitStatus, rusage *Rusage) error {
	<-ph.haveStatus
	if ph.err != nil {
		return ph.err
	}
	*status = ph.status
	if rusage != nil {
		*rusage = ph.rusage
	}
	return nil
}

func stopProcessReaper() bool {
	ForkLock.Lock()
	defer ForkLock.Unlock()
	if len(processHandles) > 0 {
		// forkExec() has been called after the last call to
		// Wait4(). Continue reaping processes.
		return false
	}
	processReaperRunning = false
	return true
}

func reapProcesses() {
	for {
		switch pid, err := blockUntilWaitable(0); err {
		case nil:
			if pid == 0 {
				// This operating system does not support
				// wait6(-1, WNOWAIT) or waitid(-1, WNOWAIT).
				// We must thus call Wait4(-1, 0). This
				// unfortunately means we can't delay waiting
				// until all calls to ProcessHandle.Hold() and
				// Unhold() have completed.
				var status WaitStatus
				var rusage Rusage
				switch pid, err := Wait4(-1, &status, 0, &rusage); err {
				case nil:
					ForkLock.Lock()
					if ph, ok := processHandles[pid]; ok {
						delete(processHandles, pid)
						ph.markTerminated(false)
						ph.status = status
						ph.rusage = rusage
						close(ph.haveStatus)
					}
					ForkLock.Unlock()
				case ECHILD:
					if stopProcessReaper() {
						return
					}
				case EINTR:
				default:
					panic(err)
				}
			} else {
				ForkLock.Lock()
				if ph, ok := processHandles[pid]; ok {
					// Process for which we have a valid
					// handle has terminated.
					delete(processHandles, pid)
					ForkLock.Unlock()
					ph.markTerminated(true)
					for {
						if _, ph.err = Wait4(pid, &ph.status, 0, &ph.rusage); ph.err != EINTR {
							break
						}
					}
					close(ph.haveStatus)
				} else {
					// Process that has been reparented
					// to us has terminated. Discard the
					// wait results.
					ForkLock.Unlock()
					var status WaitStatus
					if _, err := Wait4(pid, &status, 0, nil); err != nil && err != ECHILD && err != EINTR {
						panic(err)
					}
				}
			}
		case ECHILD:
			if stopProcessReaper() {
				return
			}
		default:
			panic(err)
		}
	}
}

var (
	processReaperInsertionLock sync.Mutex

	// These fields are locked either by acquiring ForkLock for
	// writing, or by acquiring ForkLock for reading and
	// processReaperInsertionLock.
	processReaperRunning = false
	processHandles       = map[int]*managedProcessHandle{}
)

func forkExec(argv0 string, argv []string, attr *ProcAttr, waitInBackground bool) (pid int, handle *managedProcessHandle, err error) {
	var p [2]int
	var n int
	var err1 Errno
	var wstatus WaitStatus

	if attr == nil {
		attr = &zeroProcAttr
	}
	sys := attr.Sys
	if sys == nil {
		sys = &zeroSysProcAttr
	}

	// Convert args to C form.
	argv0p, err := BytePtrFromString(argv0)
	if err != nil {
		return 0, nil, err
	}
	argvp, err := SlicePtrFromStrings(argv)
	if err != nil {
		return 0, nil, err
	}
	envvp, err := SlicePtrFromStrings(attr.Env)
	if err != nil {
		return 0, nil, err
	}

	if (runtime.GOOS == "freebsd" || runtime.GOOS == "dragonfly") && len(argv) > 0 && len(argv[0]) > len(argv0) {
		argvp[0] = argv0p
	}

	var chroot *byte
	if sys.Chroot != "" {
		chroot, err = BytePtrFromString(sys.Chroot)
		if err != nil {
			return 0, nil, err
		}
	}
	var dir *byte
	if attr.Dir != "" {
		dir, err = BytePtrFromString(attr.Dir)
		if err != nil {
			return 0, nil, err
		}
	}

	// Both Setctty and Foreground use the Ctty field,
	// but they give it slightly different meanings.
	if sys.Setctty && sys.Foreground {
		return 0, nil, errorspkg.New("both Setctty and Foreground set in SysProcAttr")
	}
	if sys.Setctty && sys.Ctty >= len(attr.Files) {
		return 0, nil, errorspkg.New("Setctty set but Ctty not valid in child")
	}

	acquireForkLock()

	// Allocate child status pipe close on exec.
	if err = forkExecPipe(p[:]); err != nil {
		releaseForkLock()
		return 0, nil, err
	}

	// Kick off child.
	pid, err1 = forkAndExecInChild(argv0p, argvp, envvp, chroot, dir, attr, sys, p[1])
	if err1 != 0 {
		Close(p[0])
		Close(p[1])
		releaseForkLock()
		return 0, nil, Errno(err1)
	}

	if waitInBackground {
		processReaperInsertionLock.Lock()
		handle = &managedProcessHandle{
			haveStatus: make(chan struct{}),
		}
		if _, ok := processHandles[pid]; ok {
			panic("Process ID has been recycled before wait status was obtained")
		}
		processHandles[pid] = handle

		if !processReaperRunning {
			processReaperRunning = true
			go reapProcesses()
		}
		processReaperInsertionLock.Unlock()
	}
	releaseForkLock()

	// Read child error status from pipe.
	Close(p[1])
	for {
		n, err = readlen(p[0], (*byte)(unsafe.Pointer(&err1)), int(unsafe.Sizeof(err1)))
		if err != EINTR {
			break
		}
	}
	Close(p[0])
	if err != nil || n != 0 {
		if n == int(unsafe.Sizeof(err1)) {
			err = Errno(err1)
		}
		if err == nil {
			err = EPIPE
		}

		if !waitInBackground {
			// Child failed; wait for it to exit, to make sure
			// the zombies don't accumulate.
			_, err1 := Wait4(pid, &wstatus, 0, nil)
			for err1 == EINTR {
				_, err1 = Wait4(pid, &wstatus, 0, nil)
			}
		}
		return 0, nil, err
	}

	// Read got EOF, so pipe closed on exec, so exec succeeded.
	return pid, handle, nil
}

var (
	// Guard the forking variable.
	forkingLock sync.Mutex
	// Number of goroutines currently forking, and thus the
	// number of goroutines holding a conceptual write lock
	// on ForkLock.
	forking int
)

// hasWaitingReaders reports whether any goroutine is waiting
// to acquire a read lock on rw. It is defined in the sync package.
func hasWaitingReaders(rw *sync.RWMutex) bool

// acquireForkLock acquires a write lock on ForkLock.
// ForkLock is exported and we've promised that during a fork
// we will call ForkLock.Lock, so that no other threads create
// new fds that are not yet close-on-exec before we fork.
// But that forces all fork calls to be serialized, which is bad.
// But we haven't promised that serialization, and it is essentially
// undetectable by other users of ForkLock, which is good.
// Avoid the serialization by ensuring that ForkLock is locked
// at the first fork and unlocked when there are no more forks.
func acquireForkLock() {
	forkingLock.Lock()
	defer forkingLock.Unlock()

	if forking == 0 {
		// There is no current write lock on ForkLock.
		ForkLock.Lock()
		forking++
		return
	}

	// ForkLock is currently locked for writing.

	if hasWaitingReaders(&ForkLock) {
		// ForkLock is locked for writing, and at least one
		// goroutine is waiting to read from it.
		// To avoid lock starvation, allow readers to proceed.
		// The simple way to do this is for us to acquire a
		// read lock. That will block us until all current
		// conceptual write locks are released.
		//
		// Note that this case is unusual on modern systems
		// with O_CLOEXEC and SOCK_CLOEXEC. On those systems
		// the standard library should never take a read
		// lock on ForkLock.

		forkingLock.Unlock()

		ForkLock.RLock()
		ForkLock.RUnlock()

		forkingLock.Lock()

		// Readers got a chance, so now take the write lock.

		if forking == 0 {
			ForkLock.Lock()
		}
	}

	forking++
}

// releaseForkLock releases the conceptual write lock on ForkLock
// acquired by acquireForkLock.
func releaseForkLock() {
	forkingLock.Lock()
	defer forkingLock.Unlock()

	if forking <= 0 {
		panic("syscall.releaseForkLock: negative count")
	}

	forking--

	if forking == 0 {
		// No more conceptual write locks.
		ForkLock.Unlock()
	}
}

// Combination of fork and exec, careful to be thread safe.
func ForkExec(argv0 string, argv []string, attr *ProcAttr) (int, error) {
	pid, _, err := forkExec(argv0, argv, attr, false)
	return pid, err
}

type ProcessHandle interface {
	Hold() (success bool)
	Unhold()
	Wait4(status *WaitStatus, rusage *Rusage) error
}

// freestandingProcessHandle is created for processes for which
// Wait4(pid) may be called directly. This is safe if the current
// process does not have PID 1, as we know no other processes will be
// reparented to it. There is thus no need to call Wait4(-1).
type freestandingProcessHandle struct {
	baseProcessHandle
	pid int
}

func (ph *freestandingProcessHandle) Wait4(status *WaitStatus, rusage *Rusage) error {
	// If we can block until Wait4 will succeed immediately, do so.
	pid, err := blockUntilWaitable(ph.pid)
	if err != nil {
		return err
	}
	if pid != 0 {
		ph.markTerminated(true)
	}
	for {
		if _, err := Wait4(ph.pid, status, 0, rusage); err == nil {
			ph.markTerminated(false)
			return nil
		} else if err != EINTR {
			return err
		}
	}
}

var (
	isPID1     bool
	isPID1Init sync.Once
)

// StartProcess wraps ForkExec for package os.
func StartProcess(argv0 string, argv []string, attr *ProcAttr) (int, ProcessHandle, error) {
	isPID1Init.Do(func() { isPID1 = Getpid() == 1 })
	if isPID1 {
		// When running as PID 1, orphan processes may be
		// reparented to us. This means that calling Wait4(pid)
		// on is not sufficient to get rid of all zombie
		// processes. Spawn a dedicated goroutine that
		// repeatedly calls Wait4(-1) and delivers the
		return forkExec(argv0, argv, attr, true)
	}

	pid, _, err := forkExec(argv0, argv, attr, false)
	if err != nil {
		return 0, nil, err
	}
	return pid, &freestandingProcessHandle{pid: pid}, nil
}

// Implemented in runtime package.
func runtime_BeforeExec()
func runtime_AfterExec()

// execveLibc is non-nil on OS using libc syscall, set to execve in exec_libc.go; this
// avoids a build dependency for other platforms.
var execveLibc func(path uintptr, argv uintptr, envp uintptr) Errno
var execveDarwin func(path *byte, argv **byte, envp **byte) error
var execveOpenBSD func(path *byte, argv **byte, envp **byte) error

// Exec invokes the execve(2) system call.
func Exec(argv0 string, argv []string, envv []string) (err error) {
	argv0p, err := BytePtrFromString(argv0)
	if err != nil {
		return err
	}
	argvp, err := SlicePtrFromStrings(argv)
	if err != nil {
		return err
	}
	envvp, err := SlicePtrFromStrings(envv)
	if err != nil {
		return err
	}
	runtime_BeforeExec()

	rlim, rlimOK := origRlimitNofile.Load().(Rlimit)
	if rlimOK && rlim.Cur != 0 {
		Setrlimit(RLIMIT_NOFILE, &rlim)
	}

	var err1 error
	if runtime.GOOS == "solaris" || runtime.GOOS == "illumos" || runtime.GOOS == "aix" {
		// RawSyscall should never be used on Solaris, illumos, or AIX.
		err1 = execveLibc(
			uintptr(unsafe.Pointer(argv0p)),
			uintptr(unsafe.Pointer(&argvp[0])),
			uintptr(unsafe.Pointer(&envvp[0])))
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "ios" {
		// Similarly on Darwin.
		err1 = execveDarwin(argv0p, &argvp[0], &envvp[0])
	} else if runtime.GOOS == "openbsd" && (runtime.GOARCH == "386" || runtime.GOARCH == "amd64" || runtime.GOARCH == "arm" || runtime.GOARCH == "arm64") {
		// Similarly on OpenBSD.
		err1 = execveOpenBSD(argv0p, &argvp[0], &envvp[0])
	} else {
		_, _, err1 = RawSyscall(SYS_EXECVE,
			uintptr(unsafe.Pointer(argv0p)),
			uintptr(unsafe.Pointer(&argvp[0])),
			uintptr(unsafe.Pointer(&envvp[0])))
	}
	runtime_AfterExec()
	return err1
}
