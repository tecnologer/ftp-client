package test

import (
	"testing"

	ftp "github.com/tecnologer/ftp-v2/src"
)

var client *ftp.Client

func init() {
	client = ftp.NewClient("54.39.115.191", ".")
}
func TestConnection(t *testing.T) {
	err := client.Connect("renechiquete@gmail.com.95750", "holamundo123.#")

	if err != nil {
		t.Fail()
	}
}

func TestFetchData(t *testing.T) {

	err := client.Connect("renechiquete@gmail.com.95750", "holamundo123.#")

	if err != nil {
		t.Fail()
	}

	entry, err := client.FetchData("/")
	if err != nil {
		t.Fail()
	}

	t.Log(entry)
}
