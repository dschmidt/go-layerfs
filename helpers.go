package layerfs

import (
	"io/fs"
)

// GetFsForDirEntry returns the source layer for a DirEntry.
func GetFsForDirEntry(d fs.DirEntry) (fs.FS, error) {
	info, err := d.Info()
	if err != nil {
		return nil, err
	}

	// Use indirection over Info() because WalkDir creates a statDirEntry
	// wrapper and there's no way we can inject our own type there
	// In contrast to that we can always provide our own fileInfo
	// and use that even from WalkDir callback.
	fileInfo, ok := info.(fileInfo)
	if !ok {
		return nil, newError("Could not assert DirEntry type", d.Name())
	}

	return fileInfo.GetFs(), nil
}
