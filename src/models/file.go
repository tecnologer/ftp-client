package models

//ElementType indentifier for the element type (directory or file)
type ElementType byte

const (
	//DirectoryType is identifier for Directory
	DirectoryType ElementType = iota
	//FileType is identifier for File
	FileType
)

//Element in the path
type Element struct {
	Name      string
	FullPath  string
	Extension string
	Type      ElementType
}
