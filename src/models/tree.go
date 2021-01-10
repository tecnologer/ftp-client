package models

import (
	"fmt"
	p "path"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/errors"
)

type Tree struct {
	*ftp.Entry
	Entries []*Tree
}

func (t *Tree) GetDirectory(path string) (*Tree, error) {
	return getBranch(t, GetPathLevels(path))
}

func (t *Tree) AddFtpEntries(path string, entries []*ftp.Entry) {
	levels := GetPathLevels(path)
	parent, _ := getBranch(t, levels)
	if parent == nil {
		parent = MkdirParent(t, levels)
	}
	parent.Entries = ParseFTPEntries(entries)
}

func (t *Tree) GetFile(path string) (*Tree, error) {
	file := p.Base(path)
	parent, err := getBranch(t, GetPathLevels(path))

	if parent == nil {
		return nil, errors.Wrap(err, "get file")
	}

	if parent.Entries == nil {
		return nil, fmt.Errorf("folder is empty")
	}

	for _, entry := range parent.Entries {
		if entry.Name == file && entry.Type == ftp.EntryTypeFile {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("file doesn't exists")
}

func getBranch(parent *Tree, levels []string) (*Tree, error) {
	if len(levels) == 0 || (len(levels) == 1 && parent.Name == levels[0]) {
		return parent, nil
	}

	if parent.Name != levels[0] {
		return nil, fmt.Errorf("folder doesn't exists")
	}

	for _, dir := range parent.Entries {
		if dir.Name == levels[1] {
			return getBranch(dir, levels[1:])
		}
	}

	return nil, fmt.Errorf("folder doesn't exists")
}

//Mkdir creates an entry type folder with the specific name
func Mkdir(name string) *Tree {
	return &Tree{
		Entry: &ftp.Entry{
			Name:   name,
			Target: name,
			Type:   ftp.EntryTypeFolder,
			Size:   0,
			Time:   time.Now().UTC(),
		},
		Entries: make([]*Tree, 0),
	}
}

//MkdirParent creates all children entries hierarchical into specific entry
func MkdirParent(parent *Tree, branches []string) *Tree {
	if len(branches) == 0 {
		return parent
	}

	if parent == nil {
		parent = Mkdir(branches[0])
		branches = branches[1:]
	}

	if parent.Entries == nil {
		parent.Entries = make([]*Tree, 0)
	}

	child := Mkdir(branches[0])
	_ = MkdirParent(child, branches[1:])
	parent.Entries = append(parent.Entries, child)
	return parent
}

//ParseFTPEntries parses []ftp.Entry to []Entry
func ParseFTPEntries(ftpEntries []*ftp.Entry) []*Tree {
	entries := make([]*Tree, len(ftpEntries))

	for i, ftpEntry := range ftpEntries {
		entries[i] = &Tree{
			Entry:   ftpEntry,
			Entries: make([]*Tree, 0),
		}
	}

	return entries
}
