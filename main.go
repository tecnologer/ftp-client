package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
)

var (
	username   string
	password   string
	host       string
	port       int
	startPath  string
	c          *ftp.ServerConn
	fileCount  int
	totalBytes int64
	reqVersion bool
	minversion string
	version    string
)

func init() {
	flag.StringVar(&username, "user", "", "(Required) username for credentials")
	flag.StringVar(&host, "host", "", "(Required) URL to the server")
	flag.StringVar(&password, "pwd", "", "password for credentials")
	flag.IntVar(&port, "port", 21, "port to connect")
	flag.StringVar(&startPath, "path", "/", "location of files in the server")
	flag.BoolVar(&reqVersion, "version", false, "returns the current version")

	flag.Parse()
}

func main() {
	if reqVersion {
		logrus.Info(version + minversion)
		return
	}
	var err error
	err = validateFlags()
	if err != nil {
		flag.PrintDefaults()
		logrus.Fatal(err)
	}

	url := fmt.Sprintf("%s:%d", host, port)

	c, err = ftp.Dial(url, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		logrus.Fatal(err)
	}

	err = c.Login(username, password)
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	defer func() {
		fmt.Printf("\n%s downloaded (%d files)  in %v", byteCountDecimal(totalBytes), fileCount, time.Since(startTime))
	}()

	if err = downloadContent(startPath); err != nil {
		logrus.Fatal(err)
	}

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
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
			logrus.Info("new folder found ", elementPath)
			if err = downloadContent(elementPath); err != nil {
				return err
			}
		}

		if element.Type != ftp.EntryTypeFile {
			continue
		}

		if err = writeFile(elementPath); err != nil {
			return err
		}
	}
	return nil
}

func writeFile(filename string) error {
	logrus.Info("downloading file ", filename)
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

	logrus.Info("file ", filename, " download sucessfully")
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
