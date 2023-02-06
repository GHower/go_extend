package taildir

import (
	"bufio"
	"gopkg.in/tomb.v1"
	"os"
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
		Filename: filename,
		file:     nil,
		reader:   nil,
		Tomb:     tomb.Tomb{},
		lk:       sync.RWMutex{},
	}
	return tail, nil
}
