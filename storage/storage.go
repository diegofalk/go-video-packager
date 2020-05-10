package storage

import (
	"io"
	"os"
)

const localSavePath string = "content/"

type UploadedContent struct {
	Name string
	Data io.ReadCloser
}

func (uc *UploadedContent) Save() error {
	file, err := os.OpenFile(localSavePath+uc.Name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	io.Copy(file, uc.Data)

	return nil
}
