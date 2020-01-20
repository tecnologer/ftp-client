package settings

//Env struct to store environment settings
type Env struct {
	NeedWait   bool `json:"need_wait"`
	ReqVersion bool
	Store      bool
}

//NewEnv returns a new instance of Env struct
func NewEnv() *Env {
	return &Env{
		NeedWait:   false,
		ReqVersion: false,
		Store:      false,
	}
}
