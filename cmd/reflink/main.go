package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gitlab.com/moreljul/reflink"
)

var flagsMap = map[string]reflink.CloneFlags{
	"mode":       reflink.PreserveMode,
	"ownership":  reflink.PreserveOwner,
	"timestamps": reflink.PreserveTimes,
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <source> <target>\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", r)
			os.Exit(1)
		}
	}()

	preserve := flag.String("preserve", "mode,ownership,timestamps", "Preserve the specified attributes (allowed values: mode, ownership, timestamps, all)")
	noPreserve := flag.String("no-preserve", "", "Don't preserve the specified attributes")
	fail := flag.Bool("fail", false, "Fail when attributes couldn't be preserved")
	unlinkOnFailure := flag.Bool("fail-unlink", false, "Unlink the target file on failure (even on attribute copy)")
	mode := flag.String("mode", "644", "Set the mode of the target file (unless -preserve=mode is used)")
	srcOffset := flag.Uint64("src-offset", 0, "Link only a subset of the source file, starting at the given offset. Must be used with -src-length")
	srcLength := flag.Uint64("src-length", 0, "Link at most the given amount of data of the source")
	destOffset := flag.Uint64("dest-offset", 0, "Offset in the destination file. Use with -src-length and/or -src-offset")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		return
	}
	source, target := args[0], args[1]

	intMode, err := parseMode(*mode)
	if err != nil {
		panic(err)
	}

	var flags reflink.CloneFlags
	if err := parseFlags(&flags, *preserve, true); err != nil {
		panic(err)
	}
	if err := parseFlags(&flags, *noPreserve, false); err != nil {
		panic(err)
	}
	if !*fail {
		flags = flags.Set(reflink.FailSafe)
	}
	if *unlinkOnFailure {
		flags = flags.Set(reflink.UnlinkOnFailure)
	}

	switch {
	case *srcLength != 0:
		err = reflink.CloneRange(source, target, flags, intMode, reflink.Offsets{
			SrcOffset:  *srcOffset,
			SrcLength:  *srcLength,
			DestOffset: *destOffset})

	case *srcOffset != 0:
		err = errors.New("-src-length must be used with -src-offset")

	case *destOffset != 0:
		err = errors.New("-dest-offset must be used with -src-length and/or -src-offset")

	default:
		err = reflink.Clone(source, target, flags, intMode)
	}

	if err != nil {
		panic(err)
	}
}

func parseMode(mode string) (os.FileMode, error) {
	intMode, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("Invalid mode: %s", mode)
	}
	return os.FileMode(intMode), nil
}

func parseFlags(flags *reflink.CloneFlags, values string, set bool) error {
	if values == "" {
		return nil
	}

	for _, name := range strings.Split(values, ",") {
		flag, ok := flagsMap[strings.ToLower(name)]
		if !ok {
			return fmt.Errorf("Unknown flag: %s", name)
		}
		if set {
			// Enable the flag
			*flags = flags.Set(flag)
		} else {
			// Disable the flag
			*flags = flags.UnSet(flag)
		}
	}

	return nil
}
