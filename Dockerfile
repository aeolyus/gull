FROM golang:alpine
LABEL maintainer="aeolyus"

WORKDIR /app
COPY . .

RUN apk add --no-cache git gcc musl-dev
RUN go get -d -v ./...
RUN go install -v ./...

VOLUME /app/data
EXPOSE 8081

CMD ["gull"]
