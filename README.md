# Go reflink library

Package reflink provides a simple library allowing to create reflinks (shallow copies) of files enabling
copy-on-write behavior.

This library currently only supports creating reflinks on Linux and requires a compatible filesystem (e.g. btrfs)

A small command line utility is also provided.

## Installation

Install using `go get`:

```
go get github.com/morelj/reflink
```

## Usage

Use the `Clone` function to create a reflink:

```go
import "reflink"

func main() {
	// Clone /source to /target
	err := reflink.Clone("/source", "/target", 0, 0644)
	if err != nil {
		panic(err)
	}

	// Clone /source to /target2, preserving source mode & owner if possible
	err = Clone("/source", "/target2", reflink.PreserveMode|reflink.PreserveOwner|reflink.FailSafe, 0644)
	if err != nil {
		panic(err)
	}
}
```

You may also only clone a subset of a file using `CloneRange`:

```go
import "reflink"

func main() {
	// Clone the first 512 bytes of /source to /target
	err := reflink.CloneRange("/source", "/target", 0, 0644, reflink.Offsets{
		SrcOffset:  0,
		SrcLength:  512,
		DestOffset: 0,
	})
	if err != nil {
		panic(err)
	}
}
```

On unsupported operating systems, `Clone` and `CloneRange` always return `reflink.ErrUnsupported`.
