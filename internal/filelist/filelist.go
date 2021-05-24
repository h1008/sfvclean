package filelist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	rxp, err := regexp.Compile(`~\d{8}-\d{6}`)
	if err != nil {
		return err
	}
	for _, file := range fl {
		if !filepath.IsAbs(file) {
			return fmt.Errorf("path is not absolute: %s", file)
		}

		if !pathContainsPart(file, ".stversions") {
			return fmt.Errorf("path is not part of .stversions folder: %s", file)
		}

		if !rxp.MatchString(file) {
			return fmt.Errorf("file name does not contain timestamp: %s", file)
		}

		fi, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("failed to call stat on file %s: %w", file, err)
		}

		if !fi.Mode().IsRegular() {
			return fmt.Errorf("not a regular file %s", file)
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
