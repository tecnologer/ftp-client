package test

import (
	"testing"

	"github.com/google/uuid"
	ftp "github.com/tecnologer/ftp-v2/src"
	"github.com/tecnologer/go-secrets"
	"github.com/tecnologer/go-secrets/config"
)

var client *ftp.Client
var ftpSecrets secrets.Secret

func init() {
	var err error
	secrets.InitWithConfig(&config.Config{BucketID: uuid.MustParse("78d50a08-c31c-4e22-9c38-6a7ada7c3649")})
	ftpSecrets, err = secrets.GetGroup("ftp")
	if err != nil {
		panic(err)
	}
	client = ftp.NewClient(ftpSecrets.GetString("host"), ".")
}

func TestConnection(t *testing.T) {
	err := client.Connect(ftpSecrets.GetString("username"), ftpSecrets.GetString("password"))

	if err != nil {
		t.Fail()
	}
}

func TestFetchData(t *testing.T) {

	err := client.Connect(ftpSecrets.GetString("username"), ftpSecrets.GetString("password"))

	if err != nil {
		t.Fail()
	}

	entry, err := client.GetEntries("/")
	if err != nil {
		t.Fail()
	}

	t.Log(entry)
}
