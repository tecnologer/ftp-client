package settings

import (
	"flag"
	"fmt"
)

var (
	username  string
	password  string
	host      string
	port      int
	startPath string
	destPath  string

	reqVersion bool
	needWait   bool
	store      bool
)

func init() {
	flag.StringVar(&username, "user", "", "(Required) username for credentials")
	flag.StringVar(&host, "host", "", "(Required) URL to the server")
	flag.StringVar(&password, "pwd", "", "password for credentials")
	flag.IntVar(&port, "port", 21, "port to connect")
	flag.StringVar(&startPath, "path", "/", "location of files in the server")
	flag.StringVar(&destPath, "dest-path", ".", "location to save the files in local")
	flag.BoolVar(&reqVersion, "version", false, "returns the current version")
	flag.BoolVar(&needWait, "wait", false, "prevents the program exit on finish process")
	flag.BoolVar(&store, "store", false, "store flags config into settings file")

	flag.Parse()
}

//Validate validates the configuration
func validateFlags() error {
	if username == "" {
		return fmt.Errorf("username is required")
	}

	if host == "" {
		return fmt.Errorf("host is required")
	}

	return nil
}
