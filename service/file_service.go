package service

import (
	"context"
	"errors"
	"io"
	"os"
	"path"
)

type FileService struct {
	cfg *Config
}

func (f *FileService) Find(ctx context.Context, fileName string) (rc io.ReadCloser, err error) {
	rc, err = f.cfg.localStorage.Seek(ctx, path.Join(csvFolder, fileName))
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return rc, nil
}
