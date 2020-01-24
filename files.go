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
func (f *filesDownload) Add(filename string) {
	if filename == "" {
		return
	}
	f.files = append(f.files, filename)
}

func (f *filesDownload) Len() int {
	return len(f.files)
}

func (f *filesDownload) HasFiles() bool {
	return f.currentIndex < f.Len()
}

func (f *filesDownload) GetNext() string {
	if !f.HasFiles() {
		return ""
	}
	file := f.files[f.currentIndex]
	f.currentIndex++
	return file
}
