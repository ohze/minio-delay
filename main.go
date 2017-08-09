package main

import (
	"github.com/minio/minio-go"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/ohze/minio-delayed-server/api"
)

func main() {
	c, err := minio.New(
		os.Getenv("MD_ENDPOINT"),
		os.Getenv("MD_KEY"),
		os.Getenv("MD_SECRET"),
		strings.ToLower(os.Getenv("MD_HTTPS")) == "true")
	if err != nil {
		log.Fatalln(err)
	}
	d := api.New(c, os.Getenv("MD_BUCKET_NAME"))
	err = http.ListenAndServe(os.Getenv("MD_PORT"), d.Handler())
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		log.Printf("listening %v", os.Getenv("MD_PORT")) //hmm. This is not printed
	}
}
