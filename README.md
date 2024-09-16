# cryptotrack

## Started command
```
chmod +x run.sh
./run.sh
```
## Cronjobs
```
docker run --rm -v /projects/cryptotrack/db.sqlite3:/srv/db.sqlite3 cryptotrack -mode=update
```
