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

type folderWorkerInput struct {
	ch     chan *remoteContent
	isBusy bool
}

var (
	remoteData *remoteContentList
)

func init() {
	remoteData = newRemoteContentList()
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

	foldersChannel := []*folderWorkerInput{}
	results := make(chan *dataError)

	for w := 1; w <= workerCount; w++ {
		folder := &folderWorkerInput{
			ch:     make(chan *remoteContent),
			isBusy: false,
		}
		foldersChannel = append(foldersChannel, folder)
		go fetchFoldersListWorker(w, folder, results, w-1)
	}

	//initial workers
	for _, fCh := range foldersChannel {
		if !remoteData.hasEntries() {
			break
		}

		fCh.ch <- remoteData.getNext()
	}

	for remoteData.hasEntries() || workerRunning(foldersChannel) {
		select {
		case result := <-results:
			if result.err != nil {
				logrus.Debugf("error donwloading file %s. Error: %v\n", result.path, result.err)
			}

			foldersChannel[result.index].ch <- remoteData.getNext()
		}
	}

	//close channels
	for _, fCh := range foldersChannel {
		close(fCh.ch)
	}

	return nil
}

func workerRunning(foldersChannel []*folderWorkerInput) bool {
	for _, fCh := range foldersChannel {
		if fCh.isBusy {
			return true
		}
	}

	return false
}

func fetchFoldersListWorker(id int, entry *folderWorkerInput, results chan<- *dataError, index int) {
	c, err := newFtpClient(config)
	if err != nil {
		logrus.WithError(err).Errorf("connecting worker #%d to FTP server", index)
		return
	}
	defer c.Quit()

	for rc := range entry.ch {
		entry.isBusy = true
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

		entry.isBusy = false
		results <- &dataError{
			err:   err,
			path:  rc.String(),
			index: index,
		}

	}
}
