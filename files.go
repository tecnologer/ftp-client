package main

type filesDownload struct {
	files        []string
	currentIndex int
}

var (
	filesToDownload *filesDownload
)

func init() {
	resetFiles()
}

func resetFiles() {
	filesToDownload = &filesDownload{
		files:        []string{},
		currentIndex: 0,
	}
}
func (f *filesDownload) add(filename string) {
	if filename == "" {
		return
	}
	f.files = append(f.files, filename)
}

func (f *filesDownload) len() int {
	return len(f.files)
}

func (f *filesDownload) hasFiles() bool {
	return f.currentIndex < f.len()
}

func (f *filesDownload) getNext() string {
	if !f.hasFiles() {
		return ""
	}
	file := f.files[f.currentIndex]
	f.currentIndex++
	return file
}
