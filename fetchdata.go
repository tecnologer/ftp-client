package main

import (
	"runtime"

	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
)

type dataError struct {
	err   error
	path  string
	index int
}

var (
	remoteData    *remoteContentList
	channelStatus map[int]bool
)

func init() {
	remoteData = newRemoteContentList()
	channelStatus = make(map[int]bool)
}

// func retriveFiles(ftpClient *ftp.ServerConn, path string) error {
// 	content, err := ftpClient.List(path)

// 	if err != nil {
// 		return err
// 	}

// 	remoteData.append(path, content...)

// 	for _, element := range content {
// 		elementPath := fmt.Sprintf("%s/%s", path, element.Name)

// 		if strings.HasSuffix(elementPath, "..") || strings.HasSuffix(elementPath, ".") {
// 			continue
// 		}

// 		if element.Type == ftp.EntryTypeFolder {
// 			// logrus.Info("new folder found ", elementPath)
// 			if err = retriveFiles(ftpClient, elementPath); err != nil {
// 				return err
// 			}
// 		}

// 		if element.Type != ftp.EntryTypeFile || config.IgnoreFile(elementPath) {
// 			continue
// 		}

// 		filesToDownload.add(elementPath)

// 		// if err = writeFile(elementPath); err != nil {
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }

func fetchDataProcess(ftpClient *ftp.ServerConn, path string) error {
	content, err := ftpClient.List(path)

	if err != nil {
		return err
	}

	remoteData.append(path, content...)

	workerCount := runtime.NumCPU()

	foldersChannel := []chan *remoteContent{}
	results := make(chan *dataError)

	for w := 1; w <= workerCount; w++ {
		folder := make(chan *remoteContent)
		foldersChannel = append(foldersChannel, folder)
		go fetchFoldersListWorker(w, folder, results, w-1)
	}

	//initial workers
	for _, ch := range foldersChannel {
		if !remoteData.hasEntries() {
			break
		}

		ch <- remoteData.getNext()
	}

	for remoteData.hasEntries() || workerRunning() {
		select {
		case result := <-results:
			if result.err != nil {
				logrus.Debugf("error donwloading file %s. Error: %v\n", result.path, result.err)
			}

			foldersChannel[result.index] <- remoteData.getNext()
		}
	}

	//close channels
	for _, ch := range foldersChannel {
		close(ch)
	}

	return nil
}

func workerRunning() bool {
	for _, isRunning := range channelStatus {
		if isRunning {
			return true
		}
	}

	return false
}

func fetchFoldersListWorker(id int, entry <-chan *remoteContent, results chan<- *dataError, index int) {
	c, err := newFtpClient(config)
	if err != nil {
		logrus.WithError(err).Errorf("connecting worker #%d to FTP server", index)
		return
	}
	defer c.Quit()

	for rc := range entry {
		channelStatus[index] = true
		if rc.entry.Type == ftp.EntryTypeFile {
			filesToDownload.add(rc.String())
		} else if rc.entry.Type == ftp.EntryTypeFolder {
			folderPath := rc.String()
			content := []*ftp.Entry{}
			content, err = c.List(folderPath)

			if err == nil {
				remoteData.append(folderPath, content...)
			}
		}

		channelStatus[index] = false
		results <- &dataError{
			err:   err,
			path:  rc.String(),
			index: index,
		}

	}
}
