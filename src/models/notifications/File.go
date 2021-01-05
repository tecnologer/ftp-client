package notifications

type FileStatus byte

const (
	Downloaded FileStatus = iota
)

//File notification struct
type File struct {
	*Notification
	Path   string
	Size   uint64
	Status FileStatus
}

//NewNotifFile creates a notification type file
func NewNotifFile(path string, size uint64, status FileStatus, metadata *Metadata) INotification {
	return &File{
		Path:         path,
		Size:         size,
		Status:       status,
		Notification: NewNotif(FileType, metadata),
	}
}
