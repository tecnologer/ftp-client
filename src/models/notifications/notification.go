package notifications

//NotificationType <-
type NotificationType string

const (
	//FileType notification type for new file
	FileType NotificationType = "File"
	//FolderType notification type for new folder
	FolderType = "Folder"
	//ErrorType notification type for errors
	ErrorType = "Error"
	//GenericType generic notification type
	GenericType = "Generic"
)

//Notification generic struct
type Notification struct {
	Type     NotificationType
	Metadata *Metadata
}

//INotification interface for notifications
type INotification interface {
	//GetType returns the type of notification
	GetType() NotificationType
	//HasError returns if has an error
	HasError() bool
	//GetMetadata returns the metadata of notification
	GetMetadata() *Metadata
	//HasMetadata returns if the notification has metadata
	HasMetadata() bool
}

//NewNotif creates a generic notification
func NewNotif(_type NotificationType, metadata *Metadata) *Notification {
	return &Notification{
		Type:     _type,
		Metadata: metadata,
	}
}

//GetType returns the type of notification
func (n *Notification) GetType() NotificationType {
	return n.Type
}

//HasError returns if the notification has an error
func (n *Notification) HasError() bool {
	return false
}

//GetMetadata returns the metadata of notification
func (n *Notification) GetMetadata() *Metadata {
	return n.Metadata
}

//HasMetadata returns if the notification has metadata
func (n *Notification) HasMetadata() bool {
	return n.Metadata != nil
}
