package bolt

import (
	"errors"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"

	"github.com/filebrowser/filebrowser/v2/rename"
)

type renameBackend struct {
	db *storm.DB
}

func (b renameBackend) Save(r *rename.Rename) error {
	return b.db.Save(r)
}

func (b renameBackend) GetByDirPath(dirPath string) ([]*rename.Rename, error) {
	var v []*rename.Rename
	err := b.db.Select(q.Eq("DirPath", dirPath)).Find(&v)
	if errors.Is(err, storm.ErrNotFound) {
		return nil, storm.ErrNotFound
	}
	return v, err
}

func (b renameBackend) GetByDirAndNewName(dirPath, newName string) (*rename.Rename, error) {
	var v rename.Rename
	err := b.db.Select(q.Eq("DirPath", dirPath), q.Eq("NewName", newName)).First(&v)
	if errors.Is(err, storm.ErrNotFound) {
		return nil, storm.ErrNotFound
	}
	return &v, err
}

func (b renameBackend) DeleteByDirPath(dirPath string) error {
	var v []*rename.Rename
	err := b.db.Select(q.Eq("DirPath", dirPath)).Find(&v)
	if errors.Is(err, storm.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, r := range v {
		if err := b.db.DeleteStruct(r); err != nil {
			return err
		}
	}
	return nil
}

func (b renameBackend) DeleteByDirAndNewName(dirPath, newName string) error {
	r, err := b.GetByDirAndNewName(dirPath, newName)
	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return nil
		}
		return err
	}
	return b.db.DeleteStruct(r)
}
