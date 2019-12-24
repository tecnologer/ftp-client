package main

import (
	"flag"
	"log"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
)

var username string
var password string

func init() {
	flag.StringVar(&username, "name", "", "username to login in the server")
	flag.StringVar(&username, "u", "", "username to login in the server")

	flag.StringVar(&password, "password", "", "password to login in the server")
	flag.StringVar(&password, "p", "", "password to login in the server")

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
