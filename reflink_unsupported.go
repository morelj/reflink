// +build !linux

package reflink

import "os"

// Clone creates a reflink (clone) of the src file to the dst file.
// flags can be used to control how the clone will be created.
// Unless the PreserveMode flag is set, the target file mode is set to mode.
func Clone(src, dst string, flags CloneFlags, mode os.FileMode) error {
	return ErrUnsupported
}

// CloneRange creates a reflink (clone) of a sub-range of the src file to the dst file.
// ofs define offsets of the source and target file.
// flags can be used to control how the clone will be created.
// Unless the PreserveMode flag is set, the target file mode is set to mode.
func CloneRange(src, dst string, flags CloneFlags, mode os.FileMode, ofs Offsets) error {
	return ErrUnsupported
}
