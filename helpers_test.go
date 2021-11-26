package layerfs

import (
	"io/fs"
	"testing"

	"github.com/psanford/memfs"
	"github.com/stretchr/testify/assert"
)

type wrappingDirEntry struct {
	d fs.DirEntry
}

func (w wrappingDirEntry) Name() string               { return w.d.Name() }
func (w wrappingDirEntry) IsDir() bool                { return w.d.IsDir() }
func (w wrappingDirEntry) Type() fs.FileMode          { return w.d.Type() }
func (w wrappingDirEntry) Info() (fs.FileInfo, error) { return w.d.Info() }

func TestGetFsForDirEntry(t *testing.T) {
	assert := assert.New(t)
	layerFs := setupTestFs(assert)
	entries, err := fs.ReadDir(layerFs, ".")
	assert.Nil(err)
	entry := entries[0]

	// test GetFsForDirEntry works when one of our own dirEntry instances is passed
	// directly to GetLayerForDirEntry
	layer, err := GetLayerForDirEntry(entry)
	assert.Nil(err)
	assert.Equal(layerFs.layers[0], layer)

	// test GetFsForDirEntry works when dirEntry is wrapped in a different type
	// e.g., WalkDir wraps the root dirEntry in a statDirEntry, so we cannot rely
	// on GetFs being present on the dirEntry itselfs
	layer, err = GetLayerForDirEntry(wrappingDirEntry{entry})
	assert.Nil(err)
	assert.Equal(layerFs.layers[0], layer)

	// test error retrieving FileInfo is propagated
	invalidDirEntry := &wrappingDirEntry{&DirEntry{
		DirEntry: &invalidDirEntryInfo{},
	}}
	_, err = GetLayerForDirEntry(invalidDirEntry)
	assert.Equal("Info failed.", err.Error())

	// test GetFsForDirEntry returns a proper error if the DirEntry
	// provides a FileInfo that cannot be asserted to be our own fileInfo
	fsys := memfs.New()
	assert.Nil(fsys.MkdirAll("dir1", 0777))
	entries, err = fs.ReadDir(fsys, ".")
	_, err = GetLayerForDirEntry(entries[0])
	assert.Equal("go-layerfs: Could not assert DirEntry type: dir1", err.Error())
}
