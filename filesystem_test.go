package kman

import (
	"os"
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func newMockFilesystem(t *testing.T, files map[string]string) afero.Fs {

	fs := afero.NewMemMapFs()
	extendMockFilesystem(t, fs, files)

	return fs
}

func extendMockFilesystem(t *testing.T, fs afero.Fs, files map[string]string) afero.Fs {

	for path, content := range files {
		require.Nil(t, afero.WriteFile(fs, path, []byte(content), os.ModePerm))
	}

	return fs
}

func snapshotFilesystem(t *testing.T, fs afero.Fs) {

	files := make(map[string]string)

	afero.Walk(fs, ".", func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() && info.Size() > 0 {

			d, err := afero.ReadFile(fs, path)
			require.Nil(t, err)

			files[path] = string(d)
		}

		return nil
	})

	snaptest.Snapshot(t, files)
}
