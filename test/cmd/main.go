package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	ftp "github.com/tecnologer/ftp-v2/src"
	"github.com/tecnologer/ftp-v2/src/models"
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

	entry, err := client.GetEntriesRecursively("/Tecnologerland")
	// entry, err = client.GetEntries("/Tecnologerland/datapacks")

	// var metadata notif.Metadata
	// for notification := range client.Notifications {
	// 	fields := logrus.Fields{
	// 		"Type":     notification.GetType(),
	// 		"HasError": notification.HasError(),
	// 	}

	// 	if notification.HasMetadata() {
	// 		metadata = *notification.GetMetadata()
	// 		for key, val := range metadata {
	// 			if key == "sizeDownloaded" {
	// 				size, _ := val.(uint64)
	// 				fields[key] = byteCountDecimal(size)
	// 			} else {
	// 				fields[key] = val
	// 			}
	// 		}

	// 		if v, ok := metadata["status"]; ok && v == "completed" {
	// 			logrus.Infof("completed in %v", metadata["duration"])
	// 			break
	// 		}
	// 	}

	// 	logrus.WithFields(fields).Info("new notification")
	// }
	// if err != nil {
	// 	panic(err)
	// }
	// logrus.WithFields(logrus.Fields{
	// 	"Name":   entry.Name,
	// 	"Target": entry.Target,
	// 	"Type":   entry.Type,
	// 	"Size":   entry.Size,
	// 	"Time":   entry.Time,
	// }).Info("entry")
	printTree(entry, 0)
	// _type := ""
	// for _, child := range entry.Entries {
	// 	if child.Type == models.EntryTypeFile {
	// 		_type = "File"
	// 	} else if child.Type == models.EntryTypeFolder {
	// 		_type = "Foder"
	// 	}
	// 	logrus.WithFields(logrus.Fields{
	// 		"Name":   child.Name,
	// 		"Target": child.Target,
	// 		"Type":   _type,
	// 		"Size":   child.Size,
	// 		"Time":   child.Time,
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

func printTree(folder *models.TreeElement, deep int) {
	if folder == nil {
		return
	}

	under := "-"
	markFolder := ""
	markFile := ""
	if deep > 0 {
		under = "|_"
		markFolder = strings.Repeat("|", deep)
		markFile = strings.Repeat("|", deep+1)

	}

	formatFolder := fmt.Sprintf("%s%%%ds %%s\n", markFolder, deep*2)
	formatFile := fmt.Sprintf("%s%%%ds %%s\n", markFile, (deep+1)*2)
	fmt.Printf(formatFolder, under, folder.Name)

	for _, entry := range folder.Entries {
		if entry.Type == models.EntryTypeFolder {
			printTree(entry, deep+1)
			continue
		}

		fmt.Printf(formatFile, "|_", entry.Name)
	}
}
