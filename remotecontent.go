package main

import (
	"fmt"

	"github.com/jlaffaye/ftp"
)

type remoteContentList struct {
	entries []*remoteContent
	current int
}

type remoteContent struct {
	entry *ftp.Entry
	path  string
}

func newRemoteContentList() *remoteContentList {
	return &remoteContentList{
		entries: []*remoteContent{},
		current: 0,
	}
}

func (rc *remoteContentList) append(path string, data ...*ftp.Entry) {
	newEntries := []*remoteContent{}
	for _, entry := range data {
		newEntry := &remoteContent{
			entry: entry,
			path:  path,
		}

		newEntries = append(newEntries, newEntry)
	}
	rc.entries = append(rc.entries, newEntries...)
}

func (rc *remoteContentList) hasEntries() bool {
	return rc.current < len(rc.entries)
}

func (rc *remoteContentList) getNext() *remoteContent {
	if !rc.hasEntries() {
		return nil
	}
	defer func() { rc.current++ }()
	return rc.entries[rc.current]
}

func (rc *remoteContent) String() string {
	return fmt.Sprintf("%s/%s", rc.path, rc.entry.Name)
}
