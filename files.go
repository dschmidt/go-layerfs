package layerfs

import (
	"errors"
	"io/fs"
)

type file struct {
	fs.File

	fs fs.FS
}

func (f *file) GetFs() fs.FS {
	return f.fs
}

type dirFile struct {
	file

	layerFs *layerFs
	name    string
}

func (f *dirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if n >= 0 {
		return nil, errors.New("layerFilesystem: Could not ReadDir because n >= 0 is not supported")
	}

	return f.layerFs.ReadDir(f.name)
}

type fileInfo struct {
	fs.FileInfo

	fs fs.FS
}

func (f *fileInfo) GetFs() fs.FS {
	return f.fs
}

type dirEntry struct {
	fs.DirEntry

	fs fs.FS
}

func (e *dirEntry) GetFs() fs.FS {
	return e.fs
}

func (e *dirEntry) Info() (fs.FileInfo, error) {
	info, err := e.DirEntry.Info()
	if err != nil {
		return nil, err
	}

	return fileInfo{
		info,
		e.fs,
	}, nil
}
