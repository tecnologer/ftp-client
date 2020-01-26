package main

import (
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/tecnologer/ftp-client/settings"
)

func newFtpClient(conf *settings.Config) (*ftp.ServerConn, error) {
	url := conf.GetURL()

	c, err := ftp.Dial(url, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	err = c.Login(config.FTP.Username, config.FTP.Password)
	if err != nil {
		return nil, err
	}

	return c, nil
}
