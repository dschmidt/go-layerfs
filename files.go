package layerfs

import (
	"io/fs"
)

type dirFile struct {
	fs.File
	fs fs.FS

	layerFs *LayerFs
	name    string
}

func (f *dirFile) GetFs() fs.FS {
	return f.fs
}

func (f *dirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if n >= 0 {
		return nil, newError("could not ReadDir because n >= 0 is not supported", f.name)
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, newError("could not ReadDir because dirFile does not point to a directory", f.name)
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
