package ftp

import (
	"fmt"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/errors"
	"github.com/tecnologer/ftp-v2/src/models"
	notif "github.com/tecnologer/ftp-v2/src/models/notifications"
)

//Client is the struct for instance of FTP client
type Client struct {
	URL           string
	Username      string
	Password      string
	Entries       *models.Entries
	PlainEntries  []string
	Timeout       time.Duration
	DestPath      string
	RootPath      string
	Notifications chan notif.INotification
	DownloadStats chan uint64

	connection   *ftp.ServerConn
	host         string
	port         int
	lock         sync.Mutex
	plainEntryCh chan string
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
		plainEntryCh:  make(chan string),
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

func (c *Client) RefreshEntriesAsync(rootPath string) {
	go func() {
		go c.refreshEntries(rootPath)

		c.PlainEntries = make([]string, 0)

		for entry := range c.plainEntryCh {
			c.lock.Lock()
			c.PlainEntries = append(c.PlainEntries, entry)
			c.lock.Unlock()
		}
	}()
}

//GetEntries updates the list of data from the server
func (c *Client) GetEntries(rootPath string) (*models.Entry, error) {
	// if c.connection == nil {
	// 	return nil, errors.Errorf("fetching data: the connection is required")
	// }

	// _, err := getEntries(c.connection, rootPath)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Fetching data: listing")
	// }

	// // parent := c.Entries.GetEntries(rootPath)

	// parentEntry := models.MkdirParent(nil, models.GetPathLevels(rootPath))
	//models.AddEntries(parentEntry, models.ParseFTPEntries(content))
	return nil, nil
}

//DownloadAsync downloads the file in the specified directory or the specific file
func (c *Client) DownloadAsync(path string, recursively bool) {
	go c.download(path, recursively)
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

//isValidFolder returns true if is not current or previous folder, this to prevent infinity recursivity
func isValidFolder(folderPath string) bool {
	return folderPath != "." && folderPath != ".."
}
