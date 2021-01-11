package test

import (
	"fmt"
	"testing"

	"github.com/jlaffaye/ftp"
	"github.com/tecnologer/ftp-v2/src/models"
)

func TestMkdir(t *testing.T) {
	output := &models.TreeElement{
		Entry: &ftp.Entry{
			Name: "/",
		},
		Entries: []*models.TreeElement{
			{
				Entry: &ftp.Entry{
					Name: "test",
				},
				Entries: []*models.TreeElement{
					{
						Entry: &ftp.Entry{
							Name: "folder1",
						},
						Entries: []*models.TreeElement{},
					},
				},
			},
		},
	}

	input := "/test/folder1/"

	levels := models.GetPathLevels(input)
	result := models.MkdirParent(nil, levels)

	if !compareTree(output, result) {
		t.Fail()
	}
}

func TestGetDirectory(t *testing.T) {
	input := &models.TreeElement{
		Entry: &ftp.Entry{
			Name: "/",
		},
		Entries: []*models.TreeElement{
			{
				Entry: &ftp.Entry{
					Name: "test",
				},
				Entries: []*models.TreeElement{
					{
						Entry: &ftp.Entry{
							Name: "folder1",
						},
						Entries: []*models.TreeElement{},
					},
					{
						Entry: &ftp.Entry{
							Name: "folder2",
						},
						Entries: []*models.TreeElement{
							{
								Entry: &ftp.Entry{
									Name: "subfolder1",
								},
								Entries: []*models.TreeElement{},
							},
						},
					},
				},
			},
		},
	}

	// output := &models.Tree{
	// 	Entry: &ftp.Entry{
	// 		Name: "folder2",
	// 	},
	// 	Entries: []*models.Tree{
	// 		{
	// 			Entry: &ftp.Entry{
	// 				Name: "subfolder1",
	// 			},
	// 			Entries: []*models.Tree{},
	// 		},
	// 	},
	// }
	output := &models.TreeElement{
		Entry: &ftp.Entry{
			Name: "subfolder1",
		},
		Entries: []*models.TreeElement{},
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

func compareTree(left, rigth *models.TreeElement) bool {
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
