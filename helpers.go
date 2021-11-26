package layerfs

import (
	"io/fs"
)

// use indirection over Info() because WalkDir creates a statDirEntry
// wrapper and there's no way we can inject our own type there.
func GetSourceFsForDirEntry(d fs.DirEntry) (fs.FS, error) {
	info, err := d.Info()
	if err != nil {
		return nil, err
	}

	fileInfo, ok := info.(fileInfo)
	if !ok {
		return nil, newError("Could not assert DirEntry type", d.Name())
	}

	return fileInfo.GetFs(), nil
}
