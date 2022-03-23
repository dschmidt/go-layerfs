package layerfs

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/psanford/memfs"
	"github.com/stretchr/testify/assert"
)

type invalidFile struct {
}

func (*invalidFile) Stat() (fs.FileInfo, error) {
	return nil, errors.New("Stat failed")
}

func (*invalidFile) Read([]byte) (int, error) {
	return 0, errors.New("Read failed")
}

func (*invalidFile) Close() error {
	return errors.New("Close failed")
}

type invalidDirEntryInfo struct {
}

func (*invalidDirEntryInfo) Name() string {
	return "invalid"
}
func (*invalidDirEntryInfo) IsDir() bool {
	return false
}
func (*invalidDirEntryInfo) Type() fs.FileMode {
	return 0
}
func (*invalidDirEntryInfo) Info() (fs.FileInfo, error) {
	return nil, errors.New("Info failed.")
}

func TestDirFileGetFs(t *testing.T) {
	assert := assert.New(t)

	fsys := memfs.New()
	dirFile := DirFile{
		name: "dir1",
		fs:   fsys,
	}
	assert.Equal(fsys, dirFile.GetFs())
}

func TestDirFileReadDir(t *testing.T) {
	assert := assert.New(t)

	fsys := memfs.New()

	dirFile := DirFile{
		File: &invalidFile{},
		fs:   fsys,
		name: "foo.txt",
	}

	// test ReadDir does not strictly positive numbers as argument
	_, err := dirFile.ReadDir(1)
	assert.Equal(err.Error(), "go-layerfs: could not ReadDir because n > 0 is not supported: foo.txt")

	// test ReadDir propagates Stat errors
	_, err = dirFile.ReadDir(-1)
	assert.Equal(err.Error(), "Stat failed")

	// test ReadDir errors out if it does not point to a dir
	assert.Nil(fsys.WriteFile("foo.txt", []byte("bar"), 0755))
	f, err := fsys.Open("foo.txt")
	assert.Nil(err)
	dirFile.File = f
	_, err = dirFile.ReadDir(-1)
	assert.Equal(err.Error(), "go-layerfs: could not ReadDir because dirFile does not point to a directory: foo.txt")
}

func TestDirEntryInfo(t *testing.T) {
	assert := assert.New(t)

	dirEntry := DirEntry{
		DirEntry: &invalidDirEntryInfo{},
	}

	_, err := dirEntry.Info()
	assert.Equal(err.Error(), "Info failed.")
}
