package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	ftp "github.com/tecnologer/ftp-v2/src"
	"github.com/tecnologer/ftp-v2/src/models/files"
)

var (
	port   = flag.Int("port", 8088, "port of server")
	client *ftp.Client
	config *ftpConfig
)

type ftpConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type msgRes struct {
	Message string `json:"message"`
}

type fileReq struct {
	Path string `json:"path"`
}

func main() {
	flag.Parse()
	config = new(ftpConfig)

	http.HandleFunc("/api/connect", connectClient)
	http.HandleFunc("/api/file/ftp", getFtpFiles)
	http.HandleFunc("/api/file/local", getLocalFiles)

	host := fmt.Sprintf(":%d", *port)
	fmt.Printf("listening on %s\n", host)
	log.Fatal(http.ListenAndServe(host, nil))
}

func getFtpFiles(w http.ResponseWriter, res *http.Request) {
	setHeaders(&w)
	var err error
	pathQ, exists := res.URL.Query()["path"]

	if !exists || len(pathQ) == 0 {
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write(newMsgRes("file path is required"))
		return
	}

	recursivelyQ, exists := res.URL.Query()["recursively"]
	recursively := !exists || len(recursivelyQ) == 0
	if len(recursivelyQ) > 0 {
		recursively, err = strconv.ParseBool(strings.ToLower(recursivelyQ[0]))
		if err != nil {
			recursively = true
		}
	}
	path := pathQ[0]
	var files *files.TreeElement
	if recursively {
		files, err = client.GetEntriesRecursively(path)
	} else {
		files, err = client.GetEntries(path)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newMsgResf("error getting files: %v", err))
		return
	}

	body, err := json.Marshal(files)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newMsgResf("error parsing to json: %v", err))
		return
	}
	w.Write(body)
}

func getLocalFiles(w http.ResponseWriter, res *http.Request) {
	setHeaders(&w)
	var err error
	pathQ, exists := res.URL.Query()["path"]

	if !exists || len(pathQ) == 0 {
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write(newMsgRes("get local files: the file path is required"))
		return
	}

	recursivelyQ, exists := res.URL.Query()["recursively"]
	recursively := !exists || len(recursivelyQ) == 0
	if len(recursivelyQ) > 0 {
		recursively, err = strconv.ParseBool(strings.ToLower(recursivelyQ[0]))
		if err != nil {
			recursively = true
		}
	}
	path := pathQ[0]
	files, err := files.ListFiles(path)
	fmt.Println(recursively)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newMsgResf("get local files: error getting local files: %v", err))
		return
	}

	body, err := json.Marshal(files)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newMsgResf("get local files: error parsing to json: %v", err))
		return
	}
	w.Write(body)
}

func connectClient(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)

	resBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write(newMsgRes("FTP configuration required"))
		return
	}

	err = json.Unmarshal(resBody, config)
	if err != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write(newMsgRes("invalid FTP configuration"))
		return
	}

	if config.Port == 0 {
		config.Port = 21
	}

	client = ftp.NewClient(config.Host)
	err = client.Connect(config.Username, config.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newMsgResf("error connecting to ftp://%s:%d: %v", config.Host, config.Port, err))
		return
	}
	w.Write(newMsgResf("connection created to ftp://%s:%d", config.Host, config.Port))
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Add("Content-Type", "application/json")
}

func newMsgRes(msg string) []byte {
	err := &msgRes{
		Message: msg,
	}

	errBody, _ := json.Marshal(err)

	return errBody
}

func newMsgResf(format string, v ...interface{}) []byte {
	return newMsgRes(fmt.Sprintf(format, v...))
}
