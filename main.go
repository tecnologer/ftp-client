package main

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/sirupsen/logrus"
	s "github.com/tecnologer/ftp-client/settings"
	"github.com/tecnologer/ftp-client/ui"
)

var (
	minversion string
	version    string
	config     *s.Config
	debug      bool
)

func init() {
	// flag.BoolVar(&debug, "-v", false, "enable debug log")
	config = s.Load()
	// flag.Parse()

	// if debug {
	logrus.SetLevel(logrus.DebugLevel)
	// }

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {

	ui.Console()
	// start()
}

func start() {
	if config.Env.ReqVersion {
		logrus.Info(version + minversion)
		return
	}

	if config.Env.NeedWait {
		//wait key input to close
		defer wait()
	}

	// defer resetFiles()

	var err error
	err = config.Validate()
	if err != nil {
		showError(err)
		return
	}

	logrus.Infof("connecting to %s", config.GetURL())
	c, err := newFtpClient(config)
	if err != nil {
		showError(err)
		return
	}

	logrus.Info("connected")

	startTime := time.Now()
	defer func() {
		msg := fmt.Sprintf("\n%s downloaded (%d files)  in %v\n", byteCountDecimal(totalBytes), fileCount, time.Since(startTime))
		fmt.Printf(msg)
		_ = beeep.Notify("Donwload Complete", msg, "")
	}()

	logrus.Info("fetching information... please wait!")
	resetFiles()

	if err = fetchDataProcess(c, config.FTP.RootPath); err != nil {
		showError(err)
		return
	}

	if err = downloadMarkedFiles(); err != nil {
		showError(err)
		return
	}

	if err := c.Quit(); err != nil {
		showError(err)
		return
	}
}

func showError(err error) {
	logrus.Error(err)
	_ = beeep.Notify("Error", fmt.Sprintf("%v", err), "")
}

func wait() {
	fmt.Println("\nPress Enter to exit...")
	fmt.Scanf("\n")
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
