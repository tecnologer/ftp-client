package files

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

//ListFiles returns a tree elements in the specific path
func ListFiles(path string, recursively bool) (*TreeElement, error) {
	return listFiles(Mkdir("/"), path, recursively)
}

//listFiles append the child to the parent in the specific path
func listFiles(parent *TreeElement, path string, recursively bool) (*TreeElement, error) {
	if parent == nil {
		return nil, fmt.Errorf("list files: parent is required")
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading the directory")
	}

	elements := make([]*TreeElement, 0)
	for _, f := range files {
		eleType := EntryTypeFile
		if f.IsDir() {
			eleType = EntryTypeFolder
		}
		child := &TreeElement{
			Name:    f.Name(),
			Type:    eleType,
			Target:  "",
			Size:    uint64(f.Size()),
			Time:    f.ModTime(),
			Entries: make([]*TreeElement, 0),
		}
		elements = append(elements, child)

		if child.Type == EntryTypeFolder && recursively {
			listFiles(child, path+"/"+child.Name, recursively)
		}
	}

	parent.AddElements(path, elements)
	if err != nil {
		return nil, errors.Wrap(err, "creating tree elements")
	}
	return parent, nil
}
