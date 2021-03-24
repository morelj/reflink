// Package reflink provides a simple library allowing to create reflinks (shallow copies) of files enabling
// copy-on-write behavior.
//
// This library currently only supports creating reflinks on Linux and requires a compatible filesystem (e.g. btrfs)
//
// A small command line utility is also provided.
package reflink
