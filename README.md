# gull
A simple URL shortener made in Go
![screenshot](https://i.imgur.com/SUVa6YA.png)

## Usage
### Docker
Optionally, build the container. Otherwise, the public one will be pulled down.
```
make image
```
Run the container.
```
docker run \
    -d \
    --name gull \
    -v /gull-data/:/data/ \
    -p 8081:8081 \
    ghcr.io/aeolyus/gull:latest
```
This will preserve any persistent data under `/gull-data/`. Change this as needed.

### From Source
```
git clone https://github.com/aeolyus/gull.git
cd ./gull
make run
```
This will create a directory `./gull/data` where persistent data will be stored.

## Acknowledgements
Inspired by [mnml](https://github.com/liyasthomas/mnmlurl/)!
