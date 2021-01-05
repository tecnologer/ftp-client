package main

import (
	"github.com/sirupsen/logrus"
	ftp "github.com/tecnologer/ftp-v2/src"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	client := ftp.NewClient("54.39.115.191", "./Downloads")
	err := client.Connect("renechiquete@gmail.com.95750", "holamundo123.#")

	if err != nil {
		panic(err)
	}
	reportCh := make(chan *ftp.Reporter)
	client.DownloadAsync("/Tecnologerland", true, reportCh)

	for report := range reportCh {
		logrus.WithFields(logrus.Fields{
			"ID":   report.ID,
			"File": report.File,
			"Msg":  report.Msg,
			"Err":  report.Err,
		}).Info("new report")
	}
	// if err != nil {
	// 	panic(err)
	// }

	// for _, entry := range newEntry.Entries {
	// 	logrus.WithFields(logrus.Fields{
	// 		"Name":   entry.Name,
	// 		"Target": entry.Target,
	// 		"Type":   entry.Type,
	// 		"Size":   entry.Size,
	// 		"Time":   entry.Time,
	// 	}).Info("entry")
	// }

}
