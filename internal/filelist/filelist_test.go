package filelist

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveLoad(t *testing.T) {
	tmp, err := os.MkdirTemp("", "go-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	fl := FileList{"test1.txt", "test2.txt"}
	flPath := filepath.Join(tmp, "filelist.json")

	err = fl.Save(flPath)
	require.NoError(t, err)

	loaded, err := Load(flPath)
	require.NoError(t, err)

	assert.Equal(t, fl, loaded)
}

func TestFileList_Verify(t *testing.T) {
	tmp, err := os.MkdirTemp("", "go-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	existingFileName := filepath.Join(tmp, ".stversions", "example~20210524-185842.txt")
	os.MkdirAll(filepath.Dir(existingFileName), 0700)
	err = os.WriteFile(existingFileName, []byte{0}, 0600)
	require.NoError(t, err)

	linkFileName := filepath.Join(tmp, ".stversions", "link~20210524-192142.txt")
	err = os.Symlink(existingFileName, linkFileName)
	require.NoError(t, err)

	dirName := filepath.Join(tmp, ".stversions", "directory~20210524-192942")
	err = os.Mkdir(dirName, 0700)
	require.NoError(t, err)

	tests := []struct {
		name    string
		fl      FileList
		wantErr bool
	}{
		{name: "Empty", fl: FileList{}},
		{name: "Success", fl: FileList{existingFileName}},
		{name: "Relative path", fl: FileList{"relative/path/tofile.txt"}, wantErr: true},
		{name: "Not in .stversions", fl: FileList{"/path/without/stversions/tofile.txt"}, wantErr: true},
		{name: "Missing timestamp", fl: FileList{"/.stversions/missing-timestamp.txt"}, wantErr: true},
		{name: "Does not exist", fl: FileList{filepath.Join(tmp, ".stversions", "does-not-exist~20210524-185842.txt")}, wantErr: true},
		{name: "Symlink", fl: FileList{linkFileName}, wantErr: true},
		{name: "Directory", fl: FileList{dirName}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fl.Verify(); (err != nil) != tt.wantErr {
				t.Errorf("FileList.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
