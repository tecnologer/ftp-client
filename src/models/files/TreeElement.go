package files

import (
	"fmt"
	p "path"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/errors"
)

type TreeElement struct {
	Name    string
	Target  string // target of symbolic link
	Type    EntryType
	Size    uint64
	Time    time.Time
	Entries []*TreeElement
}

//NewTreeElement creates the root directory
func NewTreeElement() *TreeElement {
	return Mkdir("/")
}

//GetDirectory returns the entry type folder for the specific path
func (t *TreeElement) GetDirectory(path string) (*TreeElement, error) {
	return getBranch(t, GetPathLevels(path))
}

//AddFtpEntries inserts ftp entries in the specific path
func (t *TreeElement) AddFtpEntries(path string, entries []*ftp.Entry) (parent *TreeElement, err error) {
	levels := GetPathLevels(path)
	parent, _ = getBranch(t, levels)

	//build the entries hierarchy
	if parent == nil {
		parent = MkdirParent(t, levels)
		parent, err = getBranch(t, levels)
		if err != nil {
			return nil, errors.Wrap(err, "inserting FTP entries")
		}
	}

	parent.Entries = ParseFTPEntries(entries)
	return parent, nil
}

//AddElements inserts the tree elements in the specific path
func (t *TreeElement) AddElements(path string, entries []*TreeElement) (parent *TreeElement, err error) {
	levels := GetPathLevels(path)
	parent, _ = getBranch(t, levels)

	//build the entries hierarchy
	if parent == nil {
		parent = MkdirParent(t, levels)
		parent, err = getBranch(t, levels)
		if err != nil {
			return nil, errors.Wrap(err, "inserting FTP entries")
		}
	}

	parent.Entries = entries
	return parent, nil
}

//GetFile returns the entry type file in the specific path
func (t *TreeElement) GetFile(path string) (*TreeElement, error) {
	file := p.Base(path)
	parent, err := getBranch(t, GetPathLevels(path))

	if parent == nil {
		return nil, errors.Wrap(err, "get file")
	}

	if parent.Entries == nil {
		return nil, fmt.Errorf("folder is empty")
	}

	for _, entry := range parent.Entries {
		if entry.Name == file && entry.Type == EntryTypeFile {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("file doesn't exists")
}

//CountFileRecursively returns the number of files in the path and children folders
func (t *TreeElement) CountFileRecursively(path string) (int, error) {
	levels := GetPathLevels(path)
	return countFileInBranch(t, levels)
}

func getBranch(parent *TreeElement, levels []string) (*TreeElement, error) {
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
func countFileInBranch(parent *TreeElement, levels []string) (int, error) {
	if len(levels) == 0 || (len(levels) == 1 && parent.Name == levels[0]) {
		return 0, nil
	}

	if parent.Name != levels[0] {
		return 0, fmt.Errorf("folder doesn't exists")
	}

	for _, dir := range parent.Entries {
		if dir.Name == levels[1] {
			return countFileInBranch(dir, levels[1:])
		}
	}

	return 0, fmt.Errorf("folder doesn't exists")
}

//Mkdir creates an entry type folder with the specific name
func Mkdir(name string) *TreeElement {
	return &TreeElement{
		Name:    name,
		Target:  name,
		Type:    EntryTypeFolder,
		Size:    0,
		Time:    time.Now().UTC(),
		Entries: make([]*TreeElement, 0),
	}
}

//MkdirParent creates all children entries hierarchical into specific entry
func MkdirParent(parent *TreeElement, branches []string) *TreeElement {
	if len(branches) == 0 {
		return parent
	}

	if parent == nil {
		parent = Mkdir(branches[0])
		branches = branches[1:]
	}

	if branches[0] == "/" {
		branches = branches[1:]
	}

	if parent.Entries == nil {
		parent.Entries = make([]*TreeElement, 0)
	}

	child := Mkdir(branches[0])
	_ = MkdirParent(child, branches[1:])
	parent.Entries = append(parent.Entries, child)
	return parent
}

//ParseFTPEntries parses []ftp.Entry to []Entry
func ParseFTPEntries(ftpEntries []*ftp.Entry) []*TreeElement {
	entries := make([]*TreeElement, len(ftpEntries))

	for i, ftpEntry := range ftpEntries {
		entries[i] = &TreeElement{
			Name:    ftpEntry.Name,
			Target:  ftpEntry.Target,
			Type:    EntryType(ftpEntry.Type),
			Size:    ftpEntry.Size,
			Time:    ftpEntry.Time,
			Entries: make([]*TreeElement, 0),
		}
	}

	return entries
}
