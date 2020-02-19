package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gosuri/uiprogress"
	"github.com/sirupsen/logrus"

	"github.com/jlaffaye/ftp"
)

var (
	fileCount  int
	totalBytes int64
)

type fileError struct {
	err   error
	file  string
	index int
}

func downloadMarkedFiles() error {
	count := filesToDownload.len()

	if count == 0 {
		return fmt.Errorf("no files to download")
	}

	fmt.Printf("\n >>>> found %d files <<<<\n\n", count)

	workerCount := runtime.NumCPU()
	bar := getProgressBar(count)
	uiprogress.Start()
	defer uiprogress.Stop()

	filesChannel := []chan string{}
	results := make(chan *fileError)

	logrus.Debugf("creating %d workers\n", workerCount)
	for w := 1; w <= workerCount; w++ {
		files := make(chan string)
		filesChannel = append(filesChannel, files)
		go worker(w, files, results, w-1, bar)
	}

	//initial workers
	for _, ch := range filesChannel {
		if !filesToDownload.hasFiles() {
			break
		}

		ch <- filesToDownload.getNext()
	}

	for bar.Current() < count {
		select {
		case result := <-results:
			if result.err != nil {
				logrus.Debugf("error donwloading file %s. Error: %v\n", result.file, result.err)
			}

			file := filesToDownload.getNext()
			if file == "" {
				continue
			}

			filesChannel[result.index] <- file
		}
	}

	//close channels
	for _, ch := range filesChannel {
		close(ch)
	}

	return nil
}

func worker(id int, files <-chan string, results chan<- *fileError, index int, bar *uiprogress.Bar) {
	c, err := newFtpClient(config)
	if err != nil {
		logrus.WithError(err).Errorf("connecting worker #%d to FTP server", index)
		return
	}
	defer c.Quit()

	for file := range files {
		err := writeFile(c, file)
		bar.Incr()
		results <- &fileError{
			err:   err,
			file:  file,
			index: index,
		}
	}
}

func writeFile(c *ftp.ServerConn, filename string) error {
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
	filePath := ""
	pathPattern := "%s/%s"
	if strings.HasSuffix(config.FTP.DestPath, "/") {
		pathPattern = "%s%s"
	}

	filePath = fmt.Sprintf(pathPattern, config.FTP.DestPath, filename)

	err = os.MkdirAll(filepath.Dir(filePath), 0700)
	if err != nil {
		return err
	}
	totalBytes += int64(len(buf))
	err = ioutil.WriteFile(filePath, buf, 0644)

	if err != nil {
		return err
	}

	// logrus.Info("file ", filename, " download sucessfully")
	fileCount++
	return nil
}

func getProgressBar(count int) *uiprogress.Bar {
	// create bar
	// bar := pb.New(count)

	// // refresh info every second (default 200ms)
	// // bar.SetRefreshRate(time.Second)

	// // force set io.Writer, by default it's os.Stderr
	// bar.SetWriter(os.Stdout)

	// // bar will format numbers as bytes (B, KiB, MiB, etc)
	// // bar.Set(pb.Byte, true)

	// // bar use SI bytes prefix names (B, kB) instead of IEC (B, KiB)
	// bar.Set(pb.SIBytesPrefix, true)

	// // start bar
	// bar.Start()
	return uiprogress.AddBar(count).AppendCompleted().PrependElapsed()
}
