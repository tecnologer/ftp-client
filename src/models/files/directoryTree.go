package files

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

//ListFiles returns a tree elements in the specific path
func ListFiles(path string) (*TreeElement, error) {
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

		elements = append(elements, &TreeElement{
			Name:    f.Name(),
			Type:    eleType,
			Target:  "",
			Size:    uint64(f.Size()),
			Time:    f.ModTime(),
			Entries: make([]*TreeElement, 0),
		})
	}
	root := Mkdir("/")

	parent, err := root.AddElements(path, elements)
	if err != nil {
		return nil, errors.Wrap(err, "creating tree elements")
	}
	return parent, nil
}
