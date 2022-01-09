FROM golang:alpine as build
LABEL maintainer="aeolyus"

WORKDIR /app
RUN apk add --no-cache git gcc musl-dev

COPY . .
RUN go get -d -v ./...
RUN go build -ldflags '-w -extldflags=-static' -o gull

FROM scratch

WORKDIR /

COPY --from=build /app/public /public
COPY --from=build /app/gull .

VOLUME /data
EXPOSE 8081

CMD ["./gull"]
