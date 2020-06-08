package main

import (
	"errors"
	"io"
	"os"
	//"strings"
    "net/url"
    "net/http"
    "flag"
    "log"
    "fmt"
	"goftp.io/server"
)

type ForwarderDriver struct {
	TargetURL string
	server.Perm
}

type FileInfo struct {
	os.FileInfo

	mode  os.FileMode
	owner string
	group string
}

func (f *FileInfo) Mode() os.FileMode {
	return 0
}

func (f *FileInfo) Owner() string {
	return ""
}

func (f *FileInfo) Group() string {
	return ""
}

func (driver *ForwarderDriver) Init(conn *server.Conn) {
	//driver.conn = conn
}

func (driver *ForwarderDriver) ChangeDir(path string) error {
	return nil
}

func (driver *ForwarderDriver) Stat(path string) (server.FileInfo, error) {
	return nil, errors.New("Not Implemented")
}

func (driver *ForwarderDriver) ListDir(path string, callback func(server.FileInfo) error) error {
	return errors.New("Not Implemented")
}

func (driver *ForwarderDriver) DeleteDir(path string) error {
	return errors.New("Not Implemented")
}

func (driver *ForwarderDriver) DeleteFile(path string) error {
	return errors.New("Not Implemented")
}

func (driver *ForwarderDriver) Rename(fromPath string, toPath string) error {
	return errors.New("Not Implemented")
}

func (driver *ForwarderDriver) MakeDir(path string) error {
	return errors.New("Not Implemented")
}

func (driver *ForwarderDriver) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	return 0, nil, errors.New("Not Implemented")
}

func (driver *ForwarderDriver) PutFile(destPath string, data io.Reader, appendData bool) (int64, error) {
	log.Printf("Forwarding %v", destPath)
    client := &http.Client{}
    req, err := http.NewRequest("PUT", driver.TargetURL + url.QueryEscape(destPath), data)
    resp, err := client.Do(req)

    log.Printf("Status: %v, Transferred bytes: %v, errmes: %v", resp.Status, resp.Request.ContentLength, err)
    if resp.StatusCode >= 300 {
	return 0, fmt.Errorf("server returned error status %s", resp.Status)
    }
    return resp.Request.ContentLength, err
}

type ForwarderDriverFactory struct {
	TargetURL string
	server.Perm
}

func (factory *ForwarderDriverFactory) NewDriver() (server.Driver, error) {
	return &ForwarderDriver{factory.TargetURL, factory.Perm}, nil
}


func main() {

	var (
		target = flag.String("target", "", "target url")
		user = flag.String("user", "admin", "Username for login")
		pass = flag.String("pass", "123456", "Password for login")
		port = flag.Int("port", 2121, "Port")
		host = flag.String("host", "localhost", "Host")
	)
	flag.Parse()
	if *target == "" {
		log.Fatalf("Please set a target url with -target")
	}
	factory := &ForwarderDriverFactory{TargetURL: *target, Perm: server.NewSimplePerm("user","group")}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		Auth:     &server.SimpleAuth{Name: *user, Password: *pass},
	}

	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
	log.Printf("Username %v, Password %v", *user, *pass)
	server := server.NewServer(opts)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}	
}

