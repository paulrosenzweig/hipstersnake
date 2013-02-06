#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPATH=`pwd` go build -o snake src/snake/main/main.go
coffee --compile static
cd ..
scp snake/snake root@hipstersnake.com:/tmp
scp snake/index.html root@hipstersnake.com:/tmp
scp -r snake/static/ root@hipstersnake.com:/tmp
ssh root@hipstersnake.com 'svc -d /etc/service/snake/'
ssh root@hipstersnake.com 'cat ~/snake/log.txt >> ~/full_log.txt'
ssh root@hipstersnake.com 'cp /tmp/snake ~/snake/snake'
ssh root@hipstersnake.com 'cp /tmp/index.html ~/snake/index.html'
ssh root@hipstersnake.com 'rm -rf ~/snake/static; cp -r /tmp/static/ ~/snake/'
ssh root@hipstersnake.com 'svc -u /etc/service/snake/'
