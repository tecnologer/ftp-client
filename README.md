# ftp-client

FTP Client to create backups

## Build

- Windows
  - Native: `go build -o ftpclient.exe`
  - With Makefile: `make windows`
- Linux
  - Native: `go build -o ftpclient`
  - With Makefile: `make linux`

For make binary for both OS, just use `make` or `make both`

## Usage

`./ftpclient[.exe] <-user <username>> <-host <ftp-url>> [-pwd [password]] [-port <port>] [-path <path>] [-wait] [-store]`

```txt
   -host string
        (Required) URL to the server
  -path string
        location of files in the server (default "/")
  -port int
        port to connect (default 21)
  -pwd string
        password for credentials
  -store
        store flags config into settings file
  -user string
        (Required) username for credentials
  -version
        returns the current version
  -wait
        prevents the program exit on finish process

```

## Check the version

`./ftpclient[.exe] -version`

> INFO[0000] 0.1.4.202001

## TODO

- [x] Progress bar
- [x] Settings file
- [x] Improve download process
- [ ] Encode password in settings file
- [ ] Test settings file executing from shortcut
- [ ] Improve fetching data

### Dependencies

- [FTP][1]
- [Logger][2]
- [Progress Bar][3]

[1]: github.com/jlaffaye/ftp
[2]: github.com/sirupsen/logrus
[3]: github.com/cheggaaa/pb
