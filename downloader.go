package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/jlaffaye/ftp"
)

func downloadContent(path string) error {
	content, err := c.List(path)

	if err != nil {
		return err
	}

	for _, element := range content {
		elementPath := fmt.Sprintf("%s/%s", path, element.Name)

		if strings.HasSuffix(elementPath, "..") || strings.HasSuffix(elementPath, ".") {
			continue
		}

		if element.Type == ftp.EntryTypeFolder {
			// logrus.Info("new folder found ", elementPath)
			if err = downloadContent(elementPath); err != nil {
				return err
			}
		}

		if element.Type != ftp.EntryTypeFile {
			continue
		}

		filesToDownload.Add(elementPath)

		// if err = writeFile(elementPath); err != nil {
		// 	return err
		// }
	}
	return nil
}

func downloadMarkedFiles() error {
	count := filesToDownload.Len()

	if count == 0 {
		return fmt.Errorf("no files to download")
	}

	fmt.Printf("\n >>>> found %d files <<<<\n\n", count)

	bar := pb.StartNew(count)
	defer bar.Finish()

	files := make(chan string)
	results := make(chan error)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for filesToDownload.HasFiles() {
		// if err := writeFile(filesToDownload.GetNext()); err != nil {
		// 	return err
		// }
		bar.Increment()
	}

	return nil
}

func worker(id int, files <-chan string, results chan<- error) {
	for file := range files {
		results <- writeFile(file)
	}
}

func writeFile(filename string) error {
	// logrus.Info("downloading file ", filename)
	r, err := c.Retr(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	buf, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir("."+filename), 0700)
	if err != nil {
		return err
	}
	totalBytes += int64(len(buf))
	err = ioutil.WriteFile("."+filename, buf, 0644)

	if err != nil {
		return err
	}

	// logrus.Info("file ", filename, " download sucessfully")
	fileCount++
	return nil
}
