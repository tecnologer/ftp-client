binName=ftpclient
formatVersion=+.%H%M
version=`git describe --tags`

windows:
	go build -ldflags "-X main.minversion=`date -u $(formatVersion)` -X main.version=$(version)" -o "$(binName).exe"

linux:
	go build -ldflags "-X main.minversion=`date -u $(formatVersion)` -X main.version=$(version)" -o $(binName)