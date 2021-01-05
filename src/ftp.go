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

//Client is the struct for instance of FTP client
type Client struct {
	URL           string
	Username      string
	Password      string
	Entries       *models.Entries
	Timeout       time.Duration
	DestPath      string
	RootPath      string
	host          string
	port          int
	connection    *ftp.ServerConn
	Notifications chan notif.INotification
	DownloadStats chan uint64
}

//NewClient create new instance for FTPClient
func NewClient(host, dest string) *Client {
	return &Client{
		URL:           fmt.Sprintf("%s:21", host),
		host:          host,
		Timeout:       5 * time.Second,
		Entries:       new(models.Entries),
		DestPath:      dest,
		RootPath:      "/",
		Notifications: make(chan notif.INotification),
		DownloadStats: make(chan uint64),
	}
}

//UpdatePort update the port and URL
func (c *Client) UpdatePort(port int) {
	c.port = port
	c.URL = fmt.Sprintf("%s:%d", c.host, port)
}

//Connect connects to the url using the credentials
func (c *Client) Connect(user, pwd string) (err error) {
	c.Username = user
	c.Password = pwd
	c.connection, err = c.getConnection()
	if err != nil {
		return errors.Wrap(err, "ftp connect: getting connection")
	}
	return nil
}

//FetchData updates the list of data from the server
func (c *Client) FetchData(path string) (*models.Entry, error) {
	if c.connection == nil {
		err := c.Connect(c.Username, c.Password)
		if err != nil {
			return nil, errors.Wrap(err, "fetching data: connecting")
		}
	}
	content, err := c.connection.List(path)
	if err != nil {
		return nil, errors.Wrap(err, "Fetching data: listing")
	}

	return c.updateData(path, content...), nil
}

//DownloadAsync downloads the file in the specified directory or the specific file
func (c *Client) DownloadAsync(path string, recursively bool) {
	go c.download(path, recursively)
}

func (c *Client) download(path string, recursively bool) {
	startTime := time.Now()
	go c.registerStats()
	defer c.notifyComplete(startTime)
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
				err := c.writeFile(cnn, filePath)
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

func (c *Client) downloadPath(rootPath string, recursively bool, filesCh chan<- string) error {
	entry, err := c.FetchData(rootPath)

	if err != nil {
		return errors.Wrapf(err, "downloading path %s", rootPath)
	}

	var path string
	for _, subEntry := range entry.Entries {
		path = fmt.Sprintf("%s/%s", rootPath, subEntry.Name)
		if subEntry.Type == ftp.EntryTypeFile {
			filesCh <- path
		}

		//skip the next folder if it's not recursive
		if !recursively {
			continue
		}

		if subEntry.Type == ftp.EntryTypeFolder && isValidFolder(subEntry.Name) {
			c.Notifications <- notif.NewNotifFolder(path, notif.Discovered, nil)
			err := c.downloadPath(path, recursively, filesCh)
			if err != nil {
				c.Notifications <- notif.NewNotifError(err, &notif.Metadata{"path": path, "msg": "downloading path recursively"})
			}
		}
	}

	return nil
}

func (c *Client) getConnection() (*ftp.ServerConn, error) {
	cnn, err := ftp.Dial(c.URL, ftp.DialWithTimeout(c.Timeout))
	if err != nil {
		return nil, errors.Wrap(err, "get ftp client connect: dial")
	}

	err = cnn.Login(c.Username, c.Password)
	if err != nil {
		return nil, errors.Wrap(err, "get ftp client connect: login")
	}

	return cnn, nil
}

func (c *Client) updateData(path string, data ...*ftp.Entry) *models.Entry {
	i, _ := c.Entries.GetEntries(path)

	newEntry := &models.Entry{Path: path, Entries: data}
	if i == -1 {
		*c.Entries = append(*c.Entries, newEntry)
	} else {
		(*c.Entries)[i] = newEntry
	}

	return newEntry
}

func (c *Client) notifyComplete(timestamp time.Time) {
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

func (c *Client) registerStats() {
	totalFiles := 0
	var sizeDownloaded uint64 = 0

	for downloaded := range c.DownloadStats {
		totalFiles++
		sizeDownloaded += downloaded
	}

	metadata := &notif.Metadata{
		"totalFiles":     totalFiles,
		"sizeDownloaded": sizeDownloaded,
	}

	c.Notifications <- notif.NewNotif(notif.GenericType, metadata)
}

//isValidFolder returns true if is not current or previous folder, this to prevent infinity recursivity
func isValidFolder(folderPath string) bool {
	return folderPath != "." && folderPath != ".."
}
