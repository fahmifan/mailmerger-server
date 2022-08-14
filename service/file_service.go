package service

import (
	"context"
	"errors"
	"io"
	"os"
	"path"

	"gorm.io/gorm"
)

type FileService struct {
	cfg *Config
}

func (f *FileService) Find(ctx context.Context, fileName string) (rc io.ReadCloser, err error) {
	file := File{}
	err = f.cfg.db.Take(&file, "file_name = ?", fileName).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return
	}

	rc, err = f.cfg.localStorage.Seek(ctx, path.Join(csvFolder, fileName))
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return rc, nil
}
