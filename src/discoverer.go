package ftp

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tecnologer/ftp-v2/src/models"
	notif "github.com/tecnologer/ftp-v2/src/models/notifications"
)

func (c *Client) refreshEntries(path string) {
	startTime := time.Now()
	defer c.notifyRefreshCompleted(startTime)

	workerCount := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(workerCount)

	logrus.Debugf("creating worker group with %d workers for refresh entries", workerCount)

	pathsCh := make(chan string)
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			cnn, err := c.getConnection()
			if err != nil {
				c.Notifications <- notif.NewNotifError(err, &notif.Metadata{"msg": "refreshing entries, getting connection", "workerID": id})
				return
			}
			defer cnn.Quit()

			for newPath := range pathsCh {
				c.genEntriesForPath(cnn, newPath, pathsCh)
			}
			wg.Done()
		}(i)
	}

	//c.genEntriesForPath(path, pathsCh)
	close(pathsCh)
	wg.Wait()
}

func (c *Client) genEntriesForPath(cnn *ftp.ServerConn, rootPath string, pathsCh chan<- string) error {
	entries, err := getEntries(cnn, rootPath)

	if err != nil {
		return errors.Wrapf(err, "getting entries for path %s", rootPath)
	}

	var path string
	for _, subEntry := range entries {
		path = fmt.Sprintf("%s/%s", rootPath, subEntry.Name)
		if subEntry.Type == ftp.EntryTypeFile {
			c.plainEntryCh <- path
		}

		if subEntry.Type == ftp.EntryTypeFolder && isValidFolder(subEntry.Name) {
			c.Notifications <- notif.NewNotifFolder(path, notif.Discovered, nil)
			pathsCh <- path
		}
	}
	return nil
}

func (c *Client) notifyRefreshCompleted(timestamp time.Time) {
	metadata := &notif.Metadata{
		"timestamp": timestamp,
		"duration":  time.Since(timestamp),
		"status":    "refresh completed",
	}

	c.Notifications <- notif.NewNotif(notif.GenericType, metadata)
}

func (c *Client) updateEntriesForPath(path string, data ...*ftp.Entry) *models.Entry {
	// newEntry := &models.Entry{Path: path, Entries: data}
	// c.updateEntry(newEntry)
	// return newEntry
	return nil
}

func (c *Client) updateEntry(entry *models.Entry) {
	// c.lock.Lock()
	// defer c.lock.Unlock()

	// i, _ := c.Entries.GetEntries(entry.Path)

	// if i == -1 {
	// 	*c.Entries = append(*c.Entries, entry)
	// } else {
	// 	(*c.Entries)[i] = entry
	// }
}

func getEntries(cnn *ftp.ServerConn, rootPath string) ([]*ftp.Entry, error) {
	if cnn == nil {
		return nil, errors.Errorf("fetching data: the connection is required")
	}

	return cnn.List(rootPath)
}
