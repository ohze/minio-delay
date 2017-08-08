package main

import (
	"github.com/minio/minio-go"
	"log"
	"net/http"
	"os"
	"strings"
)

const bucketName = "xenforo"

var c *minio.Client

func minioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	xr, err := newXfReq(r.Form)
	if err != nil {
		log.Print(err)
		return
	}
	if xr == nil {
		return
	}
	log.Printf("upload: %v", xr)
	//http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/
	//https://github.com/Metafused/s3-fast-upload-golang/blob/master/main.go
	go func() {
		n, err := c.FPutObject(bucketName, xr.objectName, xr.file, xr.contentType)
		if err != nil {
			log.Printf("upload fail: %v\n%v:%v", xr, n, err)
		}
	}()
}

func main() {
	var err error
	c, err = minio.New(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_KEY"),
		os.Getenv("MINIO_SECRET"),
		strings.ToLower(os.Getenv("MINIO_HTTPS")) == "true")
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/", minioHandler)
	err = http.ListenAndServe(":9004", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		log.Print("listening :9004")
	}
}
