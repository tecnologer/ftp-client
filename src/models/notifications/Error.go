package notifications

//Error notification struct
type Error struct {
	*Notification
	Err error
}

//NewNotifError creates a notification type Error
func NewNotifError(err error, metadata *Metadata) INotification {
	return &Error{
		Err:          err,
		Notification: NewNotif(ErrorType, metadata),
	}
}

//HasError returns if has an error
func (e *Error) HasError() bool {
	return e.Err != nil
}
