package main

import (
	"flag"
	"log"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
)

var (
	username string
	password string
	url      string
	port     int
)

func init() {
	flag.StringVar(&username, "user", "", "username to login in the server")
	flag.StringVar(&password, "pwd", "", "password to login in the server")
	flag.StringVar(&url, "url", "", "URL to the server")
	flag.IntVar(&port, "port", 21, "port to connect")

	flag.Parse()
}

func main() {
	if username == "" {
		logrus.Fatal("username is required")
	}

	if password == "" {
		logrus.Fatal("password is required")
	}

	c, err := ftp.Dial("ftp.example.org:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		logrus.Fatal(err)
	}

	err = c.Login(username, password)
	if err != nil {
		log.Fatal(err)
	}

	content, err := c.List("/")
	if err != nil {
		panic(err)
	}

	for _, f := range content {
		println(f)
	}

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
