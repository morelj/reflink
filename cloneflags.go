package reflink

// CloneFlags are the flags
type CloneFlags int

// Available clone flags
const (
	// Preserve file permissions
	PreserveMode CloneFlags = 1

	// Preserve file owner
	PreserveOwner CloneFlags = 1 << 1

	// Preserve file times
	PreserveTimes CloneFlags = 1 << 2

	// Don't return an error if mode / owner / times cannot be set on target file
	FailSafe CloneFlags = 1 << 6

	// Delete the target file if the operation fails
	UnlinkOnFailure CloneFlags = 1 << 7
)

// Set returns a copy of cf with f set
func (cf CloneFlags) Set(f CloneFlags) CloneFlags {
	return cf | f
}

// UnSet returns a copy of cf with f unset
func (cf CloneFlags) UnSet(f CloneFlags) CloneFlags {
	return cf & ^f
}

// IsSet returns whether f is set on cf
func (cf CloneFlags) IsSet(f CloneFlags) bool {
	return cf&f != 0
}
