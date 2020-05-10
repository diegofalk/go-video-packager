package storage

import (
	"io"
	"os"
)

const localContentPath string = "content/"
const localStreamsPath string = "stream/"

func SaveContentFile(src io.Reader, name string) error {
	file, err := os.OpenFile(localContentPath+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	io.Copy(file, src)

	return nil
}

func LoadStreamFile(dst io.Writer, name string) error {
	file, err := os.OpenFile(localStreamsPath+name, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(dst, file)

	return nil
}
