package test

import (
	"fmt"
	"testing"

	"github.com/tecnologer/ftp-v2/src/models/files"
)

func TestMkdir(t *testing.T) {
	output := &files.TreeElement{
		Name: "/",
		Entries: []*files.TreeElement{
			{
				Name: "test",
				Entries: []*files.TreeElement{
					{
						Name:    "folder1",
						Entries: []*files.TreeElement{},
					},
				},
			},
		},
	}

	input := "/test/folder1/"

	levels := files.GetPathLevels(input)
	result := files.MkdirParent(nil, levels)

	if !compareTree(output, result) {
		t.Fail()
	}
}

func TestGetDirectory(t *testing.T) {
	input := &files.TreeElement{
		Name: "/",
		Entries: []*files.TreeElement{
			{
				Name: "test",
				Entries: []*files.TreeElement{
					{
						Name:    "folder1",
						Entries: []*files.TreeElement{},
					},
					{
						Name: "folder2",
						Entries: []*files.TreeElement{
							{
								Name:    "subfolder1",
								Entries: []*files.TreeElement{},
							},
						},
					},
				},
			},
		},
	}

	// output := &models.Tree{
	// 		Name: "folder2",/ 	},
	// 	Entries: []*models.Tree{
	// 		{
	// 				Name: "subfolder1",/ 			},
	// 			Entries: []*models.Tree{},
	// 		},
	// 	},
	// }
	output := &files.TreeElement{
		Name:    "subfolder1",
		Entries: []*files.TreeElement{},
	}
	inputPath := "/test/folder2/subfolder1"

	result, err := input.GetDirectory(inputPath)

	if err != nil {
		t.Fail()
		return
	}

	if !compareTree(output, result) {
		t.Fail()
	}
}

func compareTree(left, rigth *files.TreeElement) bool {
	if left.Entries == nil && rigth.Entries != nil {
		fmt.Printf("entries nil")
		return false
	}

	if len(left.Entries) != len(rigth.Entries) {
		fmt.Printf("entries diferent size")
		return false
	}

	if len(left.Entries) > 0 {
		return compareTree(left.Entries[0], rigth.Entries[0])
	}

	return left.Name == rigth.Name
}
