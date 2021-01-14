package tools

import (
	"fmt"
	"strings"

	"github.com/tecnologer/ftp-v2/src/models/files"
)

func PrintTree(folder *files.TreeElement, deep int) {
	if folder == nil {
		return
	}

	under := "-"
	markFolder := ""
	markFile := ""
	if deep > 0 {
		under = "|_"
		markFolder = strings.Repeat("|", deep)
		markFile = strings.Repeat("|", deep+1)

	}

	formatFolder := fmt.Sprintf("%s%%%ds %%s\n", markFolder, deep*2)
	formatFile := fmt.Sprintf("%s%%%ds %%s\n", markFile, (deep+1)*2)
	fmt.Printf(formatFolder, under, folder.Name)

	for _, entry := range folder.Entries {
		if entry.Type == files.EntryTypeFolder {
			PrintTree(entry, deep+1)
			continue
		}

		fmt.Printf(formatFile, "|_", entry.Name)
	}
}
