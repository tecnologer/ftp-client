package test

import (
	"testing"

	"github.com/tecnologer/ftp-v2/src/models/files"
	"github.com/tecnologer/ftp-v2/src/models/tools"
)

func TestListFiles(t *testing.T) {
	path := "/home/tecnologer/java"

	files, err := files.ListFiles(path, false)

	if err != nil {
		t.Fail()
	}

	tools.PrintTree(files, 0)
}
