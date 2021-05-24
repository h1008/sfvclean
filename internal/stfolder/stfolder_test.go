package stfolder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createFileWithHistory(t *testing.T, path, subpathPattern string, num int) {
	for i := 1; i <= num; i++ {
		s := fmt.Sprintf("~202105%02d-121212", i)
		subpath := strings.Replace(subpathPattern, "*", s, 1)
		fn := filepath.Join(path, ".stversions", subpath)

		os.MkdirAll(filepath.Dir(fn), 0700)
		err := os.WriteFile(fn, []byte{0}, 0600)
		require.NoError(t, err)

		log.Printf("Created: %s\n", fn)
	}
}

func TestSyncThingFolder_Analyze(t *testing.T) {
	tmp, err := os.MkdirTemp("", "go-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	createFileWithHistory(t, tmp, "example/test01*.txt", 3)
	createFileWithHistory(t, tmp, "example/sub/test02*.txt", 2)

	s := &SyncThingFolder{
		stversionsPath: filepath.Join(tmp, ".stversions"),
	}

	err = s.Analyze()
	require.NoError(t, err)

	test01 := s.historyItems["example/test01.txt"]
	require.NotNil(t, test01)
	require.Len(t, test01, 3)
	assert.Equal(t, "example/test01~20210501-121212.txt", test01[0].HistoryItemPath)
	assert.Equal(t, "2021-05-01 12:12:12 +0000 UTC", test01[0].Timestamp.String())
	assert.Equal(t, int64(1), test01[0].Size)
	assert.False(t, test01[0].Filtered)

	require.FailNow(t, "TODO")
}

func Test_extractFilename(t *testing.T) {
	tests := []struct {
		name      string
		givenFile string
		wantName  string
		wantTime  string
		wantErr   bool
	}{
		{
			name:      "With Extension",
			givenFile: "some/path/to/file~20200711-205951.txt",
			wantName:  "some/path/to/file.txt",
			wantTime:  "2020-07-11 20:59:51 +0000 UTC",
		},
		{
			name:      "Without Extension",
			givenFile: "some/path/to/file~20200711-205951",
			wantName:  "some/path/to/file",
			wantTime:  "2020-07-11 20:59:51 +0000 UTC",
		},
		{
			name:      "Missing Timestamp",
			givenFile: "some/path/to/file.txt",
			wantErr:   true,
		},
		{
			name:      "Invalid Timestamp",
			givenFile: "some/path/to/file~20201400-205951",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotTime, err := extractFilename(tt.givenFile)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantName, gotName)
			assert.Equal(t, tt.wantTime, gotTime.String())
		})
	}
}
