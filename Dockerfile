FROM golang as builder

ENV GO111MODULE=on
WORKDIR /app
# copy files
COPY . .
# download dependencies
RUN go mod download
# build app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM debian AS final
WORKDIR /app
# copy app exe
COPY --from=builder /app/go-video-packager .
# copy packager exe
RUN mkdir bin
COPY bin/packager-linux bin/packager-linux
# create directories
RUN mkdir content
RUN mkdir stream

EXPOSE 8081
CMD ["/app/go-video-packager"]