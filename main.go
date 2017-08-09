package main

import (
	"github.com/minio/minio-go"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/ohze/minioxf/api"
)

func main() {
	c, err := minio.New(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_KEY"),
		os.Getenv("MINIO_SECRET"),
		strings.ToLower(os.Getenv("MINIO_HTTPS")) == "true")
	if err != nil {
		log.Fatalln(err)
	}
	xf := &api.MinioXf{Client: c}
	err = http.ListenAndServe(":9004", xf.Handler())
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		log.Print("listening :9004") //hmm. This is not printed
	}
}
