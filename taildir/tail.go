package taildir

import (
	"bufio"
	"gopkg.in/tomb.v1"
	"os"
	"path/filepath"
	"sync"
)

type Tail struct {
	Filename string

	file   *os.File
	reader *bufio.Reader

	tomb.Tomb

	lk sync.RWMutex
}

func TailF(filename string) (*Tail, error) {
	tail := &Tail{
		Filename: filepath.Clean(filename),
	}
	return tail, nil
}
