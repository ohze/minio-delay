package api

import (
	"github.com/minio/minio-go"
	"log"
	"net/http"
	"os"
	"mime"
	"path/filepath"
)

const bucketName = "xenforo"

//https://golang.org/doc/effective_go.html#embedding
type MinioXf struct {
	*minio.Client
}

//minio.Client.FPutObject with metadata
func (c * MinioXf)fPut(bucketName, objectName, filePath string, metadata map[string][]string) (n int64, err error)   {
	// Open the referenced file.
	fileReader, err := os.Open(filePath)
	// If any error fail quickly here.
	if err != nil {
		return 0, err
	}
	defer fileReader.Close()

	// Save the file stat.
	fileStat, err := fileReader.Stat()
	if err != nil {
		return 0, err
	}

	// Save the file size.
	fileSize := fileStat.Size()

	// Set contentType based on filepath extension if not given or default
	// value of "binary/octet-stream" if the extension has no associated type.
	if len(metadata["Content-Type"]) == 0 {
		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		metadata["Content-Type"] = []string{contentType}
	}

	return c.PutObjectWithSize(bucketName, objectName, fileReader, fileSize, metadata, nil)
}

//post file `file` to `object` will all metadata from PostForm (include `Content-Type`)
func (c *MinioXf) PostHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		m := r.PostForm
		file := m.Get("file")
		object := m.Get("object")
		if file == "" || object == "" {
			return
		}
		delete(m, "file")
		delete(m, "object")

		log.Printf("upload: %v -> %v", file, object)
		//http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/
		//https://github.com/Metafused/s3-fast-upload-golang/blob/master/main.go
		go func() {
			n, err := c.fPut(bucketName, object, file, m)
			if err != nil {
				log.Printf("upload fail: %v -> %v,%v\n%v:%v", file, object, m, n, err)
			}
		}()
	}
}
