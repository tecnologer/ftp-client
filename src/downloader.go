package ftp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	notif "github.com/tecnologer/ftp-v2/src/models/notifications"
)

func (c *Client) download(path, destPath string, recursively bool) {
	startTime := time.Now()

	defer c.notifyDownloadComplete(startTime)
	workerCount := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(workerCount)

	logrus.Debugf("creating worker group with %d workers", workerCount)

	filesCh := make(chan string)
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			cnn, err := c.getConnection()
			if err != nil {
				c.Notifications <- notif.NewNotifError(err, &notif.Metadata{"msg": "download file, getting connection", "workerID": id})
				return
			}
			defer cnn.Quit()

			for filePath := range filesCh {
				err := c.writeFile(cnn, filePath, destPath)
				if err != nil {
					c.Notifications <- notif.NewNotifError(err, &notif.Metadata{"msg": "download file, writting file", "workerID": id})
				}
			}
			wg.Done()
		}(i)
	}

	c.downloadPath(path, recursively, filesCh)
	close(filesCh)
	wg.Wait()
}

func (c *Client) writeFile(cnn *ftp.ServerConn, filename, destPath string) error {
	res, err := cnn.Retr(filename)
	if err != nil {
		return errors.Wrap(err, "write file: retriving file")
	}
	defer res.Close()

	buf, err := ioutil.ReadAll(res)

	if err != nil {
		return errors.Wrap(err, "write file: reading buffer")
	}

	filePath := c.getDestPath(filename, destPath)
	err = os.MkdirAll(filepath.Dir(filePath), 0700)
	if err != nil {
		return errors.Wrap(err, "write file: creating directories")
	}

	err = ioutil.WriteFile(filePath, buf, 0644)
	if err != nil {
		return errors.Wrap(err, "writing file")
	}
	c.Notifications <- notif.NewNotifFile(filePath, uint64(len(buf)), notif.Downloaded, nil)
	c.DownloadStats <- uint64(len(buf))
	return nil
}

func (c *Client) downloadPath(rootPath string, recursively bool, filesCh chan<- string) error {
	// i, entries := c.Entries.GetEntries(rootPath)

	// if i == -1 {
	// 	return errors.Errorf("not found. path %s", rootPath)
	// }

	// var path string
	// for _, subEntry := range entries {
	// 	path = fmt.Sprintf("%s/%s", rootPath, subEntry.Name)
	// 	if subEntry.Type == ftp.EntryTypeFile {
	// 		filesCh <- path
	// 	}

	// 	//skip the next folder if it's not recursive
	// 	if !recursively {
	// 		continue
	// 	}

	// 	if subEntry.Type == ftp.EntryTypeFolder && isValidFolder(subEntry.Name) {
	// 		c.Notifications <- notif.NewNotifFolder(path, notif.Discovered, nil)
	// 		err := c.downloadPath(path, recursively, filesCh)
	// 		if err != nil {
	// 			c.Notifications <- notif.NewNotifError(err, &notif.Metadata{"path": path, "msg": "downloading path recursively"})
	// 		}
	// 	}
	// }

	return nil
}

func (c *Client) getDestPath(filename, destPath string) string {
	pathPattern := "%s/%s"
	if strings.HasSuffix(destPath, "/") || strings.HasPrefix(filename, "/") {
		pathPattern = "%s%s"
	}

	filePath := fmt.Sprintf(pathPattern, destPath, filename)

	return filePath
}

func (c *Client) notifyDownloadComplete(timestamp time.Time) {
	metadata := &notif.Metadata{
		"timestamp": timestamp,
		"duration":  time.Since(timestamp),
		"status":    "completed",
	}

	for len(c.DownloadStats) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	close(c.DownloadStats)
	time.Sleep(100 * time.Millisecond)

	c.Notifications <- notif.NewNotif(notif.GenericType, metadata)
	c.DownloadStats = make(chan uint64)
}
