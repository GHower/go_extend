package watch

import "gopkg.in/tomb.v1"

type FileWatcher interface {
	BlockUntilExists(t *tomb.Tomb) error
	ChangeEvents(t *tomb.Tomb, pos int64) (*FileChanges, error)
}
