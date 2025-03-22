#!/bin/bash

set -e

gf build main.go -a amd64 -s linux -p ./temp
gf docker main.go -p -t lyy0709/auditlimit:latest
now=$(date +"%Y%m%d%H%M%S")
# 以当前时间为版本号
docker tag lyy0709/auditlimit:latest lyy0709/auditlimit:$now
docker push lyy0709/auditlimit:$now
echo "release success" $now
