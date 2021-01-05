package notifications

type FolderStatus byte

const (
	Discovered FolderStatus = iota
)

//Folder notification struct
type Folder struct {
	*Notification
	Path   string
	Status FolderStatus
}

//NewNotifFolder creates a notification type Folder
func NewNotifFolder(path string, status FolderStatus, metadata *Metadata) INotification {
	return &Folder{
		Path:         path,
		Status:       status,
		Notification: NewNotif(FolderType, metadata),
	}
}
