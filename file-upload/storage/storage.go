package storage

import (
	"github.com/ducknightii/cakes/file-upload/storage/bos"
	"io"
)

type StorageInterface interface {
	UploadFromStream(bucketName, objectName, contentType string, reader io.Reader) (string, error)
	Upload(bucketName, objectName, contentType string, body []byte) (string, error)
	MultiUpload(bucketName, objectName, contentType string, fileSize int64, reader io.Reader) (string, error)
}

var Storage StorageInterface

func Init() {
	var err error
	//Storage = &disk.Disk{BaseUrl: "localhost:18800/%s"}
	Storage, err = bos.Init("***", "***", "https://su.bcebos.com", "https://%s.su.bcebos.com")
	if err != nil {
		panic(err)
	}
}
