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
	URL        string
	Username   string
	Password   string
	Entries    *models.Entries
	Timeout    time.Duration
	DestPath   string
	RootPath   string
	host       string
	port       int
	connection *ftp.ServerConn
}

//NewClient create new instance for FTPClient
func NewClient(host, dest string) *Client {
	return &Client{
		URL:      fmt.Sprintf("%s:21", host),
		host:     host,
		Timeout:  5 * time.Second,
		Entries:  new(models.Entries),
		DestPath: dest,
		RootPath: "/",
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
func (c *Client) DownloadAsync(path string, recursively bool, reportCh chan notif.INotification) {
	go c.download(path, recursively, reportCh)
}

func (c *Client) download(path string, recursively bool, reportCh chan notif.INotification) {
	workerCount := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(workerCount)

	logrus.Debugf("creating worker group with %d workers", workerCount)

	filesCh := make(chan string)
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			cnn, err := c.getConnection()
			if err != nil {
				reportCh <- notif.NewNotifError(err, &notif.Metadata{"msg": "download file, getting connection", "workerID": id})
				return
			}
			for filePath := range filesCh {
				err := c.writeFile(cnn, filePath)
				if err != nil {
					reportCh <- notif.NewNotifError(err, &notif.Metadata{"msg": "download file, writting file", "workerID": id})
				} else {
					reportCh <- notif.NewNotifFile(filePath, 0, notif.Downloaded, nil)
				}
			}
			wg.Done()
		}(i)
	}

	c.downloadPath(path, recursively, filesCh, reportCh)
	close(filesCh)
	wg.Wait()
}

func (c *Client) downloadPath(rootPath string, recursively bool, filesCh chan<- string, reportCh chan notif.INotification) error {
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
			reportCh <- notif.NewNotifFolder(path, notif.Discovered, nil)
			err := c.downloadPath(path, recursively, filesCh, reportCh)
			if err != nil {
				reportCh <- notif.NewNotifError(err, &notif.Metadata{"path": path, "msg": "downloading path recursively"})
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

//isValidFolder returns true if is not current or previous folder, this to prevent infinity recursivity
func isValidFolder(folderPath string) bool {
	return folderPath != "." && folderPath != ".."
}
