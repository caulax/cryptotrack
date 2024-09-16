#!/bin/sh
set -e

# build
docker build -t cryptotrack .

# delete current deploy
docker rm -f cryptotrack || true

# run project
docker run -d --restart=always -v /projects/cryptotrack/db.sqlite3:/srv/db.sqlite3 --name cryptotrack -p 8080:8080 cryptotrack
