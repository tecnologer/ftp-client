package files

import (
	"bytes"
	"encoding/json"
)

// EntryType describes the different types of an Entry.
type EntryType int

// The differents types of an Entry
const (
	EntryTypeFile EntryType = iota
	EntryTypeFolder
	EntryTypeLink
)

// String returns the string representation of EntryType t.
func (t EntryType) String() string {
	return toString[t]
}

var toString = map[EntryType]string{
	EntryTypeFile:   "file",
	EntryTypeFolder: "folder",
	EntryTypeLink:   "link",
}

var toID = map[string]EntryType{
	"file":   EntryTypeFile,
	"folder": EntryTypeFolder,
	"link":   EntryTypeLink,
}

//MarshalJSON marshals the enum as a quoted json string
func (t EntryType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[t])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (t *EntryType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*t = toID[j]
	return nil
}
