package localfs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

// Storage ..
type Storage struct {
	RootDir string
}

// Save save file to local filesystem
func (s *Storage) Save(_ context.Context, dst string, reader io.Reader) error {
	dst = path.Join(s.RootDir, dst)
	dir, _ := path.Split(dst)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("unable to create dir %s : %v", dir, err)
	}

	file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("unable to create file %s : %v", dst, err)
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("unable to write file %s : %v", dst, err)
	}

	return nil
}

func (s *Storage) Seek(_ context.Context, dst string) (io.ReadCloser, error) {
	dst = path.Join(s.RootDir, dst)
	return os.Open(dst)
}
