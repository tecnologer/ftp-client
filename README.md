# ftp-client

FTP Client to create backups

## Build

- Windows
  - Native: `go build -o ftpclient.exe`
  - With Makefile: `make windows`
- Linux
  - Native: `go build -o ftpclient`
  - With Makefile: `make linux`


## Usage

`ftpclient[.exe] <-user <username>> <-host <ftp-url>> [-pwd [password]] [-port <port>] [-path <path>]`

```txt
 -host string
        (Required) URL to the server
  -path string
        location of files in the server (default "/")
  -port int
        port to connect (default 21)
  -pwd string
        password for credentials
  -user string
        (Required) username for credentials

```
