FROM golang

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN mkdir -p /app/content
RUN mkdir -p /app/stream

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 8081
CMD ["/app/go-video-packager"]