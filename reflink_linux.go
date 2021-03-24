package reflink

import (
	"os"
	"strconv"
	"unsafe"

	"golang.org/x/sys/unix"
)

const ficlone = 0x40049409
const ficlonerange = 0x4020940d

type fileCloneRange struct {
	srcFd      int64
	srcOffset  uint64
	srcLength  uint64
	destOffset uint64
}

// Clone creates a reflink (clone) of the src file to the dst file.
// flags can be used to control how the clone will be created.
// Unless the PreserveMode flag is set, the target file mode is set to mode.
func Clone(src, dst string, flags CloneFlags, mode os.FileMode) error {
	return clone(src, dst, flags, mode, func(srcFd, dstFd int) error {
		return unix.IoctlSetInt(dstFd, ficlone, srcFd)
	})
}

// CloneRange creates a reflink (clone) of a sub-range of the src file to the dst file.
// ofs define offsets of the source and target file.
// flags can be used to control how the clone will be created.
// Unless the PreserveMode flag is set, the target file mode is set to mode.
func CloneRange(src, dst string, flags CloneFlags, mode os.FileMode, ofs Offsets) error {
	return clone(src, dst, flags, mode, func(srcFd, dstFd int) error {
		arg := fileCloneRange{
			srcFd:      int64(srcFd),
			srcOffset:  ofs.SrcOffset,
			srcLength:  ofs.SrcLength,
			destOffset: ofs.DestOffset}
		_, _, err := unix.Syscall(unix.SYS_IOCTL, uintptr(dstFd), ficlonerange, uintptr(unsafe.Pointer(&arg)))
		if err != 0 {
			return err
		}
		return nil
	})
}

func clone(src, dst string, flags CloneFlags, mode os.FileMode, ioctl func(srcFd, dstFd int) error) (err error) {
	// When true, force the removal of the target on error
	forceRemove := false

	defer func() {
		// Cleanup
		if err != nil && (flags&UnlinkOnFailure != 0 || forceRemove) {
			os.Remove(dst)
		}
	}()

	// Open source file for reading
	srcFd, err := unix.Open(src, unix.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer unix.Close(srcFd)

	// Create target file
	dstFd, err := unix.Open(dst, unix.O_WRONLY|unix.O_CREAT, unix.S_IRUSR|unix.S_IWUSR)
	if err != nil {
		return
	}
	defer func() {
		cerr := unix.Close(dstFd)
		if cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Clone
	err = ioctl(srcFd, dstFd)
	if err != nil {
		forceRemove = true
		return
	}

	if flags&PreserveMode != 0 || flags&PreserveOwner != 0 || flags&PreserveTimes != 0 {
		err = cloneStat(srcFd, dstFd, flags)
		if err != nil {
			return
		}
	}

	if flags&PreserveMode == 0 {
		// Set the provided mode
		err = failSafeErr(unix.Fchmod(dstFd, uint32(mode)), flags)
	}

	return
}

func cloneStat(srcFd, dstFd int, flags CloneFlags) error {
	// Stat source file to get mode, owner and times
	var stat unix.Stat_t
	if statErr := unix.Fstat(srcFd, &stat); statErr != nil {
		return failSafeErr(statErr, flags)
	}

	if flags&PreserveMode != 0 {
		err := failSafeErr(unix.Fchmod(dstFd, stat.Mode), flags)
		if err != nil {
			return err
		}
	}

	if flags&PreserveOwner != 0 {
		err := failSafeErr(unix.Fchown(dstFd, int(stat.Uid), int(stat.Gid)), flags)
		if err != nil {
			return err
		}
	}
	if flags&PreserveMode != 0 {
		err := failSafeErr(unix.Fchmod(dstFd, stat.Mode), flags)
		if err != nil {
			return err
		}
	}
	if flags&PreserveTimes != 0 {
		// Mimic the behavior of unix.Futimes at the nanosecond level
		err := failSafeErr(unix.UtimesNano("/proc/self/fd/"+strconv.FormatInt(int64(dstFd), 10), []unix.Timespec{stat.Atim, stat.Mtim}), flags)
		if err != nil {
			return err
		}
	}

	return nil
}

func failSafeErr(err error, flags CloneFlags) error {
	if flags&FailSafe == 0 {
		return err
	}
	return nil
}
