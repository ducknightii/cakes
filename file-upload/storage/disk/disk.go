package disk

import (
	"fmt"
	"io"
	"os"

	"github.com/ducknightii/cakes/file-upload/storage/types"
)

type Disk struct {
	BaseUrl string
}

const savePath = "./upload"

func (d *Disk) Upload(bucketName, objectName, contentType string, body []byte) (string, error) {
	// store uploaded file into local path
	localFile := savePath + "/" + objectName
	out, err := os.Create(localFile)
	if err != nil {
		fmt.Printf("failed to open the file %s for writing", localFile)
		return "", err
	}
	defer out.Close()

	_, err = out.Write(body)
	if err != nil {
		fmt.Printf("wire file err:%s\n", err)
		return "", err
	}
	return fmt.Sprintf(d.BaseUrl, bucketName) + "/" + objectName, nil
}

func (d *Disk) MultiUpload(bucketName, objectName, contentType string, fileSize int64, reader types.Reader) (string, error) {

	return d.UploadFromStream(bucketName, objectName, contentType, reader)
}

func (d *Disk) UploadFromStream(bucketName, objectName, contentType string, reader types.Reader) (string, error) {
	// store uploaded file into local path
	localFile := savePath + "/" + objectName
	out, err := os.Create(localFile)
	if err != nil {
		fmt.Printf("failed to open the file %s for writing", localFile)
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		fmt.Printf("copy file err:%s\n", err)
		return "", err
	}
	return fmt.Sprintf(d.BaseUrl, bucketName) + "/" + objectName, nil
}
