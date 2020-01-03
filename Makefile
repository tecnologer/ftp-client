binName=ftpclient
formatVersion=+.%Y%m
version=`git describe --tags`

both:
	make windows
	make linux

windows:
	GOOS=windows go build -ldflags "-X main.minversion=`date -u $(formatVersion)` -X main.version=$(version)" -o "$(binName).exe"

linux:
	GOOS=linux go build -ldflags "-X main.minversion=`date -u $(formatVersion)` -X main.version=$(version)" -o $(binName)
	