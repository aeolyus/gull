FROM golang:alpine as build
LABEL maintainer="aeolyus"

WORKDIR /app
COPY . .

RUN apk add --no-cache git gcc musl-dev
RUN go get -d -v ./...
RUN go build -o gull

FROM alpine

WORKDIR /

COPY --from=build /app/public /public
COPY --from=build /app/gull .

VOLUME /data
EXPOSE 8081

CMD ["./gull"]
