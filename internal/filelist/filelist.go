package filelist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sfvclean/internal/stfolder"
)

type FileList []string

func Load(path string) (FileList, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var list FileList
	err = json.Unmarshal(content, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (fl FileList) Save(path string) error {
	bytes, err := json.MarshalIndent(fl, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0640)
}

func (fl FileList) Verify() error {
	for _, file := range fl {
		if !filepath.IsAbs(file) {
			return fmt.Errorf("path %s is not absolute", file)
		}

		if !pathContainsPart(file, ".stversions") {
			return fmt.Errorf("path %s is not part of .stversions folder", file)
		}

		if !stfolder.TimestampRegex.MatchString(file) {
			return fmt.Errorf("file name %s does not contain timestamp", file)
		}

		fi, err := os.Lstat(file)
		if err != nil {
			return fmt.Errorf("failed to call stat on file %s: %w", file, err)
		}

		if !fi.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", file)
		}

	}
	return nil
}

func pathContainsPart(path, part string) bool {
	for _, pathpart := range strings.Split(filepath.ToSlash(path), "/") {
		if pathpart == part {
			return true
		}
	}
	return false
}
