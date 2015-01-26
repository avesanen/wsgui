#!/bin/bash

docker build -t wsgui .
docker stop wsgui
docker rm wsgui
docker run -v /var/lib/wsgui:/go/src/github.com/avesanen/wsgui/db -p 127.0.0.1:5002:8000 -d --name wsgui wsgui
