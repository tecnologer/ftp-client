binName=ftpclient
formatVersion=+.%Y%m
version=`git describe --tags`

all:
	make windows
	make linux
	make darwin

windows:
	GOOS=windows go build -ldflags "-X main.minversion=`date -u $(formatVersion)` -X main.version=$(version)" -o "$(binName).exe"

linux:
	GOOS=linux go build -ldflags "-X main.minversion=`date -u $(formatVersion)` -X main.version=$(version)" -o linux-$(binName)

darwin:
	GOOS=darwin go build -ldflags "-X main.minversion=`date -u $(formatVersion)` -X main.version=$(version)" -o darwin-$(binName)
	