package main

import (
	"github.com/minio/minio-go"
	"log"
	"net/http"
	"os"
	"strings"
)

const bucketName = "xf"

var c *minio.Client

func minioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	contentType := r.Form.Get("content-type")
	file := r.Form.Get("file")
	root := r.Form.Get("root")
	if contentType == "" || file == "" || root == "" {
		return
	}
	if root[len(root)-1] != '/' {
		root += "/"
	}
	if !strings.HasPrefix(file, root) {
		return
	}
	objectName := file[len(root):]
	//http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/
	//https://github.com/Metafused/s3-fast-upload-golang/blob/master/main.go
	go c.FPutObject(bucketName, objectName, file, contentType)
	log.Printf("uploaded: %v, %v, %v", root, objectName, contentType)
}

func main() {
	var err error
	c, err = minio.New(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_KEY"),
		os.Getenv("MINIO_SECRET"),
		true)
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/", minioHandler)
	err = http.ListenAndServe(":9004", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
