package reflink

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing requires that an environment variable REFLINK_TEST_DIR_OK points to a directory on a compatible filesystem
// with write permission

func TestReflinkSuccess(t *testing.T) {
	dir := os.Getenv("REFLINK_TEST_DIR_OK")
	if dir == "" {
		t.Skip("No test OK directory (REFLINK_TEST_DIR_OK)")
	}
	if runtime.GOOS != "linux" {
		t.Skip("Unsupported OS")
	}

	assert := assert.New(t)

	// Create temporary file
	src := path.Join(dir, fmt.Sprintf("reflink-test-src-%x", rand.Int()))
	dst := path.Join(dir, fmt.Sprintf("reflink-test-dst-%x", rand.Int()))
	create(src)
	defer func() {
		os.Remove(src)
		os.Remove(dst)
	}()

	err := Clone(src, dst, PreserveMode, 0600)
	assert.NoError(err)
	info, err := os.Stat(dst)
	assert.NoError(err)
	assert.Equal(info.ModTime(), info.ModTime())
}

func create(name string) {
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	buf := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		_, err := file.Write(buf)
		if err != nil {
			panic(err)
		}
	}
}

func ExampleClone() {
	// Clone /source to /target
	err := Clone("/source", "/target", 0, 0644)
	if err != nil {
		panic(err)
	}

	// Clone /source to /target2, preserving source mode & owner if possible
	err = Clone("/source", "/target2", PreserveMode|PreserveOwner|FailSafe, 0644)
	if err != nil {
		panic(err)
	}
}

func ExampleCloneRange() {
	// Clone the first 512 bytes of /source to /target
	err := CloneRange("/source", "/target", 0, 0644, Offsets{
		SrcOffset:  0,
		SrcLength:  512,
		DestOffset: 0,
	})
	if err != nil {
		panic(err)
	}
}
