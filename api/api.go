package api

import (
	"github.com/minio/minio-go"
	"log"
	"net/http"
	"os"
	"mime"
	"path/filepath"
)

//https://golang.org/doc/effective_go.html#embedding
type DelayedApi struct {
	bucketName string
	*minio.Client
}

func New(c *minio.Client, bucketName string) *DelayedApi {
	return &DelayedApi{bucketName: bucketName, Client: c}
}

//minio.Client.FPutObject with metadata
func (d *DelayedApi)fPut(objectName, filePath string, metadata map[string][]string) (n int64, err error)   {
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

	return d.PutObjectWithSize(d.bucketName, objectName, fileReader, fileSize, metadata, nil)
}

//post file `file` to `object` will all metadata from PostForm (include `Content-Type`)
func (d *DelayedApi) Handler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			d.handlePost(w, r)
		} else if r.Method == "DELETE" {
			d.handleDelete(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	}
}

func (d *DelayedApi) handlePost(w http.ResponseWriter, r *http.Request) {
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
		n, err := d.fPut(object, file, m)
		if err != nil {
			log.Printf("upload fail: %v -> %v,%v\n%v:%v", file, object, m, n, err)
		}
	}()
}

func (d *DelayedApi) handleDelete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	objects := r.Form["objects"]
	log.Printf("delete %v. Form: %v", objects, r.Form)
	if len(objects) == 0 {
		return
	}
	go func() {
		for _, object := range objects {
			if err := d.RemoveObject(d.bucketName, object); err != nil {
				log.Printf("delete fail: %v: %v", objects, err)
			}
		}

		//hmm. Don't know why minio.Client.RemoveObjects (delete multiple object) is failed. We use RemoveObject for now
		//n := len(objects)
		//if n == 1 {
		//	if err := d.RemoveObject(bucketName, objects[0]); err != nil {
		//		log.Printf("delete fail: %v: %v", objects, err)
		//	}
		//} else {
		//	objectsCh := make(chan string, n)
		//	//defer close(objectsCh)
		//	for _, object := range objects {
		//		objectsCh <- object
		//	}
		//	errCh := d.RemoveObjects(bucketName, objectsCh)
		//	for e := range errCh {
		//		log.Printf("delete fail: %v: %v", e.ObjectName, e.Err)
		//	}
		//}
	}()
}
