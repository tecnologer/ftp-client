package main

import (
	"github.com/sirupsen/logrus"
	ftp "github.com/tecnologer/ftp-v2/src"
	notif "github.com/tecnologer/ftp-v2/src/models/notifications"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	client := ftp.NewClient("54.39.115.191", "./Downloads")
	err := client.Connect("renechiquete@gmail.com.95750", "holamundo123.#")

	if err != nil {
		panic(err)
	}
	reportCh := make(chan notif.INotification)
	client.DownloadAsync("/Tecnologerland", true, reportCh)

	for report := range reportCh {
		fields := logrus.Fields{
			"Type":     report.GetType(),
			"HasError": report.HasError(),
		}

		if report.HasMetadata() {
			for key, val := range *report.GetMetadata() {
				fields[key] = val
			}
		}

		logrus.WithFields(fields).Info("new notification")
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
