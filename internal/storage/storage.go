package storage

import "errors"

type Storage interface {
	Open() (err error)
	Load() (err error)
	Save() (err error)
}

var (
	ErrStorageOpen = errors.New("storage: error opening storage")
	ErrStorageLoad = errors.New("storage: error loading storage")
	ErrStorageSave = errors.New("storage: error saving storage")
)
