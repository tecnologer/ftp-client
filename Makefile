binName=ftpclient

windows:
	go build -o "$(binName).exe"

linux:
	go build -o $(binName)