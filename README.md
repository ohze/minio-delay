# minio-delayed-server
[![Build Status](https://travis-ci.org/ohze/minio-delayed-server.svg?branch=master)](https://travis-ci.org/ohze/minio-delayed-server)

This is a simple golang http server doing a simple thing:
When receiving a POST/ DELETE request, it asynchronously run the
operation to a provided minio server and return immediatetly.

## Usage example
```sh
docker run -e MD_ENDPOINT=lb.minio:9009 -e MD_HTTPS=false
    -e MD_KEY=minio -e MD_SECRET=minio123 \
    -e MD_BUCKET_NAME=xenforo \
    -e MD_PORT=":9004" \
    -v /var/www/html:/var/www/html:ro sandinh/minio-delayed-server

curl -X POST -F 'Content-Type=image/jpeg' -F 'file=/var/www/html/data/x.jpg' \
    -F 'object=data/y.jpg' \
    -F 'Content-Disposition=inline; filename="x.jpg"' \
    -F 'other-metadata=hello' \
    http://localhost:9004/

curl -v http://lb.minio:9009/xenforo/data/y.jpg

curl -X POST -F 'file=/var/www/html/data/y2.zip' -F 'object=data/y2.zip'

curl -X DELETE http://localhost:9004/?objects=data/y.jpg&objects=data/y2.zip
```

### Licence
This software is licensed under the Apache 2 license:
http://www.apache.org/licenses/LICENSE-2.0

Copyright 2017 Sân Đình (https://sandinh.com)
