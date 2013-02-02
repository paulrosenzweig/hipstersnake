#!/bin/bash

ssh root@hipstersnake.com 'svc -d /etc/service/snake/'
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPATH=`pwd` go build -o snake src/snake/main/main.go
cd ..
scp -r snake root@hipstersnake.com:/root
ssh root@hipstersnake.com 'svc -u /etc/service/snake/'