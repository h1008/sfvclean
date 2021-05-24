package app

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/olekukonko/tablewriter"

	"sfvclean/internal/filelist"
	"sfvclean/internal/stfolder"
	"sfvclean/internal/utils"
)

const ResultFile = "/tmp/sfvclean.json"

func Analyze(path string, dur utils.Duration, verbose bool) error {
	since := time.Now().AddDate(-dur.Years, -dur.Months, -dur.Days)
	folders, err := stfolder.Search(path)
	if err != nil {
		return nil
	}

	var fileList filelist.FileList
	for _, folder := range folders {
		err = folder.Analyze()
		if err != nil {
			return err
		}

		files, err := folder.ApplyDateFilter(since)
		if err != nil {
			return err
		}

		fileList = append(fileList, files...)
	}

	sort.Strings(fileList)

	if verbose {
		for _, file := range fileList {
			log.Println(file)
		}
	}

	err = fileList.Save(ResultFile)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Folder", "Total # Files", "Total Size", "# Files to delete", "File Size to delete"})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})

	var total stfolder.Stats
	for _, folder := range folders {
		stats := folder.Stats()

		total.TotalNum += stats.TotalNum
		total.TotalSize += stats.TotalSize
		total.FilteredNum += stats.FilteredNum
		total.FilteredSize += stats.FilteredSize

		table.Append(FormatRow(folder.Name(), stats))
	}

	if len(folders) > 1 {
		table.SetFooter(FormatRow("TOTAL", total))
		table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	}

	table.Render()
	return nil
}

func FormatRow(name string, stats stfolder.Stats) []string {
	return []string{
		name, fmt.Sprint(stats.TotalNum), utils.FormatFileSize(stats.TotalSize), fmt.Sprint(stats.FilteredNum), utils.FormatFileSize(stats.FilteredSize),
	}

}
