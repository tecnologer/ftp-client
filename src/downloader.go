package ftp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/errors"
	notif "github.com/tecnologer/ftp-v2/src/models/notifications"
)

func (c *Client) writeFile(cnn *ftp.ServerConn, filename string) error {
	res, err := cnn.Retr(filename)
	if err != nil {
		return errors.Wrap(err, "write file: retriving file")
	}
	defer res.Close()

	buf, err := ioutil.ReadAll(res)

	if err != nil {
		return errors.Wrap(err, "write file: reading buffer")
	}

	filePath := c.getDestPath(filename)
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

func (c *Client) getDestPath(filename string) string {
	pathPattern := "%s/%s"
	if strings.HasSuffix(c.DestPath, "/") || strings.HasPrefix(filename, "/") {
		pathPattern = "%s%s"
	}

	filePath := fmt.Sprintf(pathPattern, c.DestPath, filename)

	return filePath
}
