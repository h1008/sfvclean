package stfolder

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"sfvclean/internal/filelist"
	"sort"
	"strings"
	"time"
)

type SyncThingFolder struct {
	stversionsPath string
	historyItems   map[string][]HistoryItem
}

type HistoryItem struct {
	HistoryItemPath string
	Size            int64
	Timestamp       time.Time
	OriginalDeleted bool
	Filtered        bool
}

type Stats struct {
	TotalSize    int64
	TotalNum     int64
	FilteredSize int64
	FilteredNum  int64
}

func Search(root string) ([]*SyncThingFolder, error) {
	folders := make([]*SyncThingFolder, 0)
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == ".stversions" {
			folders = append(folders, &SyncThingFolder{
				stversionsPath: path,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return folders, nil
}

func (s *SyncThingFolder) Name() string {
	return filepath.Dir(s.stversionsPath)
}

func (s *SyncThingFolder) Analyze() error {
	existingFiles, err := analyzeOriginalFolder(s.Name())
	if err != nil {
		return fmt.Errorf("cannot analyze original folder: %w", err)
	}

	results, err := analyzeStversionsFolder(s.stversionsPath, existingFiles)
	if err != nil {
		return fmt.Errorf("cannot analyze stversions folder: %w", err)
	}

	s.historyItems = results

	return nil
}

func (s *SyncThingFolder) ApplyDateFilter(before time.Time) (filelist.FileList, error) {
	var fileList filelist.FileList
	for _, historyItems := range s.historyItems {
		for i := range historyItems {
			item := &historyItems[i]

			if item.Timestamp.Before(before) {
				item.Filtered = true

				itemPath := filepath.Join(s.stversionsPath, item.HistoryItemPath)
				absItemPath, err := filepath.Abs(itemPath)
				if err != nil {
					return nil, err
				}

				fileList = append(fileList, absItemPath)
			} else {
				item.Filtered = false
			}
		}
	}
	sort.Strings(fileList)
	return fileList, nil
}

func (s *SyncThingFolder) Stats() Stats {
	var stats Stats
	for _, v := range s.historyItems {
		for _, w := range v {
			stats.TotalSize += w.Size
			stats.TotalNum++

			if w.Filtered {
				stats.FilteredSize += w.Size
				stats.FilteredNum++
			}
		}
	}
	return stats
}

func analyzeOriginalFolder(folder string) (map[string]struct{}, error) {
	items := make(map[string]struct{})

	err := filepath.Walk(folder, func(path string, info fs.FileInfo, _ error) error {
		dir := filepath.Base(path)
		if dir == ".stversions" || dir == ".stfolder" {
			return filepath.SkipDir
		}
		if dir == ".stignore" {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		p, err := filepath.Rel(folder, path)
		if err != nil {
			return err
		}

		items[p] = struct{}{}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

func analyzeStversionsFolder(folder string, existingFiles map[string]struct{}) (map[string][]HistoryItem, error) {
	items := make(map[string][]HistoryItem)
	err := filepath.Walk(folder, func(path string, info fs.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		p, err := filepath.Rel(folder, path)
		if err != nil {
			return err
		}

		org, t, err := extractFilename(p)
		if err != nil {
			return err
		}

		_, ok := existingFiles[org]

		h := HistoryItem{
			HistoryItemPath: p,
			Size:            info.Size(),
			Timestamp:       t,
			OriginalDeleted: !ok,
		}

		items[org] = append(items[org], h)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

func extractFilename(file string) (string, time.Time, error) {
	rxp, err := regexp.Compile(`~(\d{8}-\d{6})`)
	if err != nil {
		return "", time.Time{}, err
	}

	matches := rxp.FindStringSubmatch(file)
	if matches != nil {
		t, err := time.ParseInLocation("20060102-150405", matches[1], time.Local)
		if err != nil {
			return "", time.Time{}, err
		}

		original := strings.Replace(file, matches[0], "", 1)
		return original, t, nil
	}

	return "", time.Time{}, fmt.Errorf("no match found: %v", file)
}
