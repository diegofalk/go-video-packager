package storage

import (
	"os"
)

const localSavePath string = "content/"

type UploadedContent struct {
	Name string
	Data []byte
}

func (uc *UploadedContent) Save() error {
	file, err := os.OpenFile(localSavePath+uc.Name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(uc.Data)

	return nil
}
