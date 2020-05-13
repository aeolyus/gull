# gull
A simple URL shortener made in Go

## Usage
### Docker
Pull the image and run.
```
docker run -d --name gull -v /gull-data/:/app/data/ -p 8081:8081 aeolyus/gull:latest
```
This will preserve any persistent data under `/gull-data/`. Change this as needed.

### From Source
```
git clone https://github.com/aeolyus/gull.git
cd ./gull
go get -d -v ./...
go run server.go
```
This will create a directory `./gull/data` where persistent data will be stored.
