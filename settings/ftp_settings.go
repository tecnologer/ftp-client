package settings

//FTP settings for File Transfer Protocol (FTP)
type FTP struct {
	Username string `json:"username"`
	Password string `json:"pwd"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	RootPath string `json:"root_path"`
	DestPath string `json:"dest_path"`
}

//NewFTP creates new instance for FTP settings
func NewFTP(username, host string) *FTP {
	return &FTP{
		Username: username,
		Host:     host,
		Password: "",
		Port:     21,
		RootPath: "/",
	}
}
