package models

import "github.com/jlaffaye/ftp"

//Entries is the group of entries in specific path
type Entries []*Entry

//Entry for element and its elements
type Entry struct {
	Path    string
	Entries []*ftp.Entry
}

//GetEntries returns the list of entries in the specific path
func (e *Entries) GetEntries(path string) (int, []*ftp.Entry) {
	for index, entry := range *e {
		if entry.Path == path {
			return index, entry.Entries
		}
	}

	return -1, make([]*ftp.Entry, 0)
}
