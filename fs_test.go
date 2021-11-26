package layerfs

import (
	"io/fs"
	"io/ioutil"
	"testing"

	"github.com/psanford/memfs"
	"github.com/stretchr/testify/assert"
)

func assertContent(assert *assert.Assertions, filesystem fs.FS, fileName string, content []byte) {
	fileContent, err := fs.ReadFile(filesystem, fileName)
	assert.Nil(err)
	assert.Equal(content, fileContent)
}

func setupTestFs(assert *assert.Assertions) *LayerFs {
	upperFs := memfs.New()
	assert.Nil(upperFs.WriteFile("f1.txt", []byte("foo"), 0755))
	assert.Nil(upperFs.WriteFile("f2.txt", []byte("foo"), 0755))
	assert.Nil(upperFs.MkdirAll("dir1", 0777))
	assert.Nil(upperFs.WriteFile("dir1/f11.txt", []byte("foo"), 0755))
	assert.Nil(upperFs.WriteFile("dir1/f12.txt", []byte("foo"), 0755))

	lowerFs := memfs.New()
	assert.Nil(lowerFs.WriteFile("f2.txt", []byte("bar"), 0755))
	assert.Nil(lowerFs.WriteFile("f3.txt", []byte("bar"), 0755))
	assert.Nil(lowerFs.MkdirAll("dir1", 0777))
	assert.Nil(lowerFs.WriteFile("dir1/f12.txt", []byte("bar"), 0755))
	assert.Nil(lowerFs.WriteFile("dir1/f13.txt", []byte("bar"), 0755))

	return New(upperFs, lowerFs)
}

func TestLayerFsOpen(t *testing.T) {
	assert := assert.New(t)
	layerFs := setupTestFs(assert)
	assertContent(assert, layerFs, "f1.txt", []byte("foo"))
	assertContent(assert, layerFs, "f2.txt", []byte("foo"))
	assertContent(assert, layerFs, "f3.txt", []byte("bar"))

	assertContent(assert, layerFs, "dir1/f11.txt", []byte("foo"))
	assertContent(assert, layerFs, "dir1/f12.txt", []byte("foo"))
	assertContent(assert, layerFs, "dir1/f13.txt", []byte("bar"))

	// test errors in upper layers are skipped
	f, err := layerFs.Open("f3.txt")
	assert.Nil(err)
	c, err := ioutil.ReadAll(f)
	assert.Nil(err)
	assert.Equal([]byte("bar"), c)

	// test proper error is returned when no layers succeed
	_, err = layerFs.Open("f4.txt")
	assert.Equal("go-layerfs: could not Open: f4.txt", err.Error())
}

func TestLayerFsReadFile(t *testing.T) {
	assert := assert.New(t)
	layerFs := setupTestFs(assert)

	// test proper error is returned when no layers succeed
	_, err := layerFs.ReadFile("f4.txt")
	assert.Equal("go-layerfs: could not ReadFile: f4.txt", err.Error())
}

func TestLayerFsReadDir(t *testing.T) {
	assert := assert.New(t)
	layerFs := setupTestFs(assert)
	_, err := fs.ReadDir(layerFs, ".")
	assert.Nil(err)

	// TODO: assert entries are correct
	// for _, e := range entries {
	// 	fmt.Printf("entry: %#v\n", e.Name())
	// }

	// test proper error is returned when no layers succeed
	_, err = fs.ReadDir(layerFs, "dir4")
	assert.Equal("go-layerfs: could not ReadDir: dir4", err.Error())
}

func TestLayerFsStat(t *testing.T) {
	assert := assert.New(t)
	layerFs := setupTestFs(assert)
	info, _ := layerFs.Stat(".")

	assert.IsType(fileInfo{}, info)

	// test errors in upper layers are skipped
	info, err := layerFs.Stat("f3.txt")
	assert.Nil(err)

	// test proper error is returned when no layers succeed
	_, err = layerFs.Stat("f4.txt")
	assert.Equal("go-layerfs: could not Stat: f4.txt", err.Error())
}

func TestLayerFsWalkDir(t *testing.T) {
	assert := assert.New(t)
	layerFs := setupTestFs(assert)
	assert.Nil(fs.WalkDir(layerFs, ".", func(path string, d fs.DirEntry, err error) error {
		assert.Nil(err)

		sourceFs, err := GetLayerForDirEntry(d)
		assert.Nil(err)
		if d.IsDir() {
			return nil
		}

		// FIXME: assert content is correct
		_, err = fs.ReadFile(sourceFs, path)
		assert.Nil(err)
		return nil
	}))
}

func TestLayerFsReadDirFile(t *testing.T) {
	assert := assert.New(t)
	layerFs := setupTestFs(assert)

	file, err := layerFs.Open(".")
	assert.Nil(err)
	readDirFile, ok := file.(fs.ReadDirFile)
	assert.True(ok)
	_, err = readDirFile.ReadDir(-1)
	assert.Nil(err)
	// FIXME: add actual assert
	// for _, e := range entries {
	// 	// fmt.Printf("entry: %#v\n", e)
	// }
}
