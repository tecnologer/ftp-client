package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	ftp "github.com/tecnologer/ftp-v2/src"
	notif "github.com/tecnologer/ftp-v2/src/models/notifications"
	"github.com/tecnologer/go-secrets"
	"github.com/tecnologer/go-secrets/config"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	secrets.InitWithConfig(&config.Config{})
	ftpSecrets, err := secrets.GetGroup("ftp")
	if err != nil {
		panic(err)
	}

	client := ftp.NewClient(ftpSecrets.GetString("host"), "./Downloads")
	err = client.Connect(ftpSecrets.GetString("username"), ftpSecrets.GetString("password"))

	if err != nil {
		panic(err)
	}

	//client.DownloadAsync("/Tecnologerland/data", true)

	entry, err := client.GetEntries("/Tecnologerland/data")

	var metadata notif.Metadata
	for notification := range client.Notifications {
		fields := logrus.Fields{
			"Type":     notification.GetType(),
			"HasError": notification.HasError(),
		}

		if notification.HasMetadata() {
			metadata = *notification.GetMetadata()
			for key, val := range metadata {
				if key == "sizeDownloaded" {
					size, _ := val.(uint64)
					fields[key] = byteCountDecimal(size)
				} else {
					fields[key] = val
				}
			}

			if v, ok := metadata["status"]; ok && v == "completed" {
				logrus.Infof("completed in %v", metadata["duration"])
				break
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

func byteCountDecimal(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
