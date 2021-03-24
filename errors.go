package reflink

type internalError string

func (e internalError) Error() string {
	return string(e)
}

func (e internalError) String() string {
	return string(e)
}

// ErrUnsupported indicated that the operation is not supported
const ErrUnsupported = internalError("reflinks are not supported on this platform")
