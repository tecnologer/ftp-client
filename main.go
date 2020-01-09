package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/gen2brain/beeep"
	"github.com/cheggaaa/pb"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
)

var (
	username        string
	password        string
	host            string
	port            int
	startPath       string
	c               *ftp.ServerConn
	fileCount       int
	totalBytes      int64
	reqVersion      bool
	needWait      bool
	minversion      string
	version         string
	filesToDownload []string
)

func init() {
	flag.StringVar(&username, "user", "", "(Required) username for credentials")
	flag.StringVar(&host, "host", "", "(Required) URL to the server")
	flag.StringVar(&password, "pwd", "", "password for credentials")
	flag.IntVar(&port, "port", 21, "port to connect")
	flag.StringVar(&startPath, "path", "/", "location of files in the server")
	flag.BoolVar(&reqVersion, "version", false, "returns the current version")
	flag.BoolVar(&needWait, "wait", false, "prevents the program exit on finish process")

	flag.Parse()
}

func main() {
	if reqVersion {
		logrus.Info(version + minversion)
		return
	}

	if needWait {
		//wait key input to close
		defer wait()
	}
	
	filesToDownload = make([]string, 0, 2)

	var err error
	err = validateFlags()
	if err != nil {
		flag.PrintDefaults()
		showError(err)
	}

	url := fmt.Sprintf("%s:%d", host, port)

	logrus.Infof("connecting to %s", url)
	c, err = ftp.Dial(url, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		showError(err)
		return
	}

	err = c.Login(username, password)
	if err != nil {
		showError(err)
		return
	}

	logrus.Info("connected")

	startTime := time.Now()
	defer func() {
		msg := fmt.Sprintf("\n%s downloaded (%d files)  in %v", byteCountDecimal(totalBytes), fileCount, time.Since(startTime))
		fmt.Printf(msg)
		_ = beeep.Notify("Donwload Complete", msg, "")
	}()

	logrus.Info("fetching information... please wait!")
	if err = downloadContent(startPath); err != nil {
		showError(err)
		return
	}

	if err = downloadMarkedFiles(); err != nil {
		showError(err)
		return
	}

	if err := c.Quit(); err != nil {
		showError(err)
		return
	}
}

func showError(err error) {
	logrus.Error(err)
	_ = beeep.Notify("Error", fmt.Sprintf("%v", err), "")
}

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

		markFileToDownload(elementPath)

		// if err = writeFile(elementPath); err != nil {
		// 	return err
		// }
	}
	return nil
}

func markFileToDownload(filename string) {
	if filename == "" {
		return
	}
	filesToDownload = append(filesToDownload, filename)
}

func downloadMarkedFiles() error {
	count := len(filesToDownload)

	if count == 0 {
		return fmt.Errorf("no files to download")
	}

	fmt.Printf("\n >>>> found %d files <<<<\n\n", count)

	bar := pb.StartNew(count)
	defer bar.Finish()

	for _, file := range filesToDownload {
		if err := writeFile(file); err != nil {
			return err
		}
		bar.Increment()
	}

	return nil
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

func validateFlags() error {
	if username == "" {
		return fmt.Errorf("username is required")
	}

	if host == "" {
		return fmt.Errorf("host is required")
	}

	return nil
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func wait(){
	fmt.Println("\nPress Enter to exit...")
	fmt.Scanf("\n")
}