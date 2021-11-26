package layerfs

import (
	"io/fs"
)

func New(layers ...fs.FS) *layerFs {
	return &layerFs{layers: layers}
}

type layerFs struct {
	layers []fs.FS
}

func (fsys *layerFs) Open(name string) (fs.File, error) {
	for _, layer := range fsys.layers {
		f, err := layer.Open(name)
		if err != nil {
			continue
		}

		return &dirFile{
			f,
			layer,
			fsys,
			name,
		}, nil
	}

	return nil, newError("could not Open", name)
}

func (fsys *layerFs) ReadFile(name string) ([]byte, error) {
	for _, layer := range fsys.layers {
		file, err := fs.ReadFile(layer, name)
		if err != nil {
			continue
		}

		return file, nil
	}

	return nil, newError("could not ReadFile", name)
}

func (fsys *layerFs) ReadDir(name string) ([]fs.DirEntry, error) {
	entryMap := make(map[string]bool)
	entries := make([]fs.DirEntry, 0)
	errorLayers := 0
	for _, layer := range fsys.layers {
		layerEntries, err := fs.ReadDir(layer, name)
		if err != nil {
			errorLayers++
			continue
		}
		for _, layerEntry := range layerEntries {
			_, ok := entryMap[layerEntry.Name()]
			if ok {
				continue
			}
			entryMap[layerEntry.Name()] = true
			lFsDirEntry := &dirEntry{
				layerEntry,
				layer,
			}
			entries = append(entries, lFsDirEntry)
		}
	}

	if errorLayers == len(fsys.layers) {
		return nil, newError("could not ReadDir", name)
	}

	return entries, nil
}

func (fsys *layerFs) Stat(name string) (fs.FileInfo, error) {
	for _, layer := range fsys.layers {
		fi, err := fs.Stat(layer, name)
		if err != nil {
			continue
		}

		return fileInfo{
			fi,
			layer,
		}, nil
	}

	return nil, newError("could not Stat", name)
}
