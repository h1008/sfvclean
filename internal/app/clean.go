package app

import (
	"fmt"
	"log"
	"os"

	"sfvclean/internal/filelist"
)

func Clean(force bool) error {
	fl, err := filelist.Load(ResultFile)
	if err != nil {
		return err
	}

	err = fl.Verify()
	if err != nil {
		return err
	}

	errors := 0
	for _, file := range fl {
		log.Printf("Deleting file %s\n", file)

		if force {
			if err := os.Remove(file); err != nil {
				log.Printf("FAILED to delete file %s: %v", file, err)
				errors++
			}
		}
	}

	if !force {
		log.Println("Dry run: No file has been deleted! Use option -force to actually delete the files.")
	}

	if errors > 0 {
		return fmt.Errorf("failed to delete %d files", errors)
	}

	return nil
}
