package main

import (
	"net/url"
	"strings"
	"errors"
	"fmt"
)

type XfReq struct {
	contentType string //ex, image/jpeg
	root string //ex, /var/www/html/data/
	file string //ex, /var/www/html/data/avatars/l/0/1.jpg
	prefix string //ex, data/

	objectName string //computed, ex, data/avatars/l/0/1.jpg
}

func (r *XfReq)String() string {
	return fmt.Sprintf("Req: %v, %v, %v", r.file, r.objectName, r.contentType)
}

func newXfReq(form url.Values) (*XfReq, error)  {
	r := new(XfReq)

	r.contentType = form.Get("content-type")
	r.file = form.Get("file")
	r.root = form.Get("root")
	r.prefix = form.Get("prefix")
	if r.contentType == "" || r.file == "" || r.root == "" || r.prefix == "" {
		return nil, nil
	}
	if r.root[len(r.root)-1] != '/' {
		r.root += "/"
	}
	if !strings.HasPrefix(r.file, r.root) {
		return nil, errors.New(fmt.Sprintf("invalid file: %v not in %v", r.file, r.root))
	}
	r.objectName = r.prefix + r.file[len(r.root):]

	return r, nil
}
