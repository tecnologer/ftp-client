package models

import (
	"strings"

	"github.com/jlaffaye/ftp"
)

//Entries is the group of entries in specific path
type Entries []*Entry

//Entry for element and its elements
type Entry struct {
	*ftp.Entry
	Entries []*Entry
}

//GetEntries returns the list of entries in the specific path
func (e *Entries) GetEntries(path string) []*Entry {
	if len(*e) == 0 {
		return nil
	}

	return getEntries(GetPathLevels(path), (*e)[0].Entries)
}

func getEntries(branches []string, entries []*Entry) []*Entry {
	if len(branches) == 0 {
		return entries
	}
	var entry *Entry
	for _, subEntries := range entries {
		if subEntries.Name == branches[0] {
			entry = subEntries
		}
	}

	if entry == nil {
		return nil
	}

	return getEntries(branches[1:], entry.Entries)
}

//GetPathLevels returns array with the branches name of the path tree
func GetPathLevels(path string) []string {
	branches := []string{}
	if !strings.HasPrefix(path, "/") {
		branches = append(branches, "/")
	}

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	branches = append(branches, strings.Split(path, "/")...)
	branches[0] = "/"
	return branches
}

// //Mkdir creates an entry type folder with the specific name
// func Mkdir(name string) *Entry {
// 	return &Entry{
// 		Entry: &ftp.Entry{
// 			Name:   name,
// 			Target: name,
// 			Type:   ftp.EntryTypeFolder,
// 			Size:   0,
// 			Time:   time.Now().UTC(),
// 		},
// 		Entries: make([]*Entry, 0),
// 	}
// }

// //MkdirParent creates all children entries hierarchical into specific entry
// func MkdirParent(parent *Entry, branches []string) *Entry {
// 	if len(branches) == 0 {
// 		return parent
// 	}

// 	if parent == nil {
// 		parent = Mkdir(branches[0])
// 		branches = branches[1:]
// 	}

// 	if parent.Entries == nil {
// 		parent.Entries = make([]*Entry, 0)
// 	}

// 	child := Mkdir(branches[0])

// 	parent.Entries = append(parent.Entries, MkdirParent(child, branches[1:]))
// 	return parent
// }
