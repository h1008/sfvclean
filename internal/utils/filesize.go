package utils

import "fmt"

const (
	KB int64 = 1024
	MB       = KB * KB
	GB       = MB * KB
	TB       = GB * KB
	PB       = TB * KB
	EB       = PB * KB
)

func FormatFileSize(size int64) string {
	if size < KB {
		return fmt.Sprintf("%d B", size)
	} else if size < MB {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	} else if size < GB {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	} else if size < TB {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	} else if size < PB {
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	} else if size < EB {
		return fmt.Sprintf("%.2f PB", float64(size)/float64(PB))
	} else {
		return fmt.Sprintf("%.2f EB", float64(size)/float64(EB))
	}
}
