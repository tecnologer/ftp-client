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
	username  string
	password  string
	host      string
	port      int
	startPath string
	c         *ftp.ServerConn
	fileCount int
)

func init() {
	flag.StringVar(&username, "user", "", "username to login in the server")
	flag.StringVar(&password, "pwd", "", "password to login in the server")
	flag.StringVar(&host, "host", "", "URL to the server")
	flag.IntVar(&port, "port", 21, "port to connect")
	flag.StringVar(&startPath, "path", "/69518", "location of files in the server")

	flag.Parse()
}

func main() {
	var err error
	err = validateFlags()
	if err != nil {
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

	// err = writeFile("/69518/plugins/Modern-LWC-2.1.5.jar")
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	startTime := time.Now()
	defer func() {
		fmt.Printf("\n%d downloaded files in %v", fileCount, time.Since(startTime))
	}()
	if err = downloadContent(startPath); err != nil {
		logrus.Fatal(err)
	}
	// content, err := c.List("/69518")
	// if err != nil {
	// 	panic(err)
	// }

	// for _, f := range content {
	// 	println(f.Name)
	// }

	// r, err := c.Retr("test-file.txt")
	// if err != nil {
	// 	panic(err)
	// }

	// buf, err := ioutil.ReadAll(r)
	// println(string(buf))

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

	if password == "" {
		return fmt.Errorf("password is required")
	}

	if host == "" {
		return fmt.Errorf("host is required")
	}

	return nil

}
