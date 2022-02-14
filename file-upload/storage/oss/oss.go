package oss

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/ducknightii/cakes/file-upload/storage/types"
	"io"
	"sync"
	"time"
)

var (
	once sync.Once
	ins  = new(client)
)

type client struct {
	client     *oss.Client
	baseUrl    string
	partSize   int64
	maxPartNum int64
}

func Init(ak, sk, endpoint, baseUrl string) (c *client, err error) {
	once.Do(func() {
		ins.baseUrl = baseUrl
		ins.partSize = 20 * 1024 * 1024 // 20Mi
		ins.maxPartNum = 100            // <= 10000
		ins.client, err = oss.New(endpoint, ak, sk)
		if err != nil {
			return
		}
	})

	return ins, err
}

func (c *client) Upload(bucketName, objectName, contentType string, body []byte) (string, error) {
	cur := time.Now()
	// todo 初始化成 bucket 结构体
	ossBucket, err := c.client.Bucket(bucketName)
	if err != nil {
		fmt.Printf("oss.Bucket err: %s\n", err)
		return "", err
	}
	opts := []oss.Option{
		oss.ContentType(contentType),
		oss.ObjectACL(oss.ACLPublicRead),
		oss.StorageClass(oss.StorageIA),
	}
	err = ossBucket.PutObject(objectName, bytes.NewReader(body), opts...)
	fmt.Printf("oss.Upload cost: %dms\n", time.Now().Sub(cur)/time.Millisecond)
	if err != nil {
		fmt.Println("oss-err", fmt.Sprintf("bucket:%s object:%s err:%s", bucketName, objectName, err))
	}

	return fmt.Sprintf(c.baseUrl, bucketName) + "/" + objectName, err
}

func (c *client) UploadFromStream(bucketName, objectName, contentType string, reader types.Reader) (string, error) {
	cur := time.Now()
	// todo 初始化成 bucket 结构体
	ossBucket, err := c.client.Bucket(bucketName)
	if err != nil {
		fmt.Printf("oss.Bucket err: %s\n", err)
		return "", err
	}
	opts := []oss.Option{
		oss.ContentType(contentType),
		oss.ObjectACL(oss.ACLPublicRead),
		oss.StorageClass(oss.StorageIA),
	}
	err = ossBucket.PutObject(objectName, reader, opts...)
	fmt.Printf("oss.UploadFromStream cost: %dms\n", time.Now().Sub(cur)/time.Millisecond)
	if err != nil {
		fmt.Println("oss-err", fmt.Sprintf("bucket:%s object:%s err:%s", bucketName, objectName, err))
	}

	return fmt.Sprintf(c.baseUrl, bucketName) + "/" + objectName, err
}

type bucket struct {
	bucket *oss.Bucket
}

type partRes struct {
	uploadPart oss.UploadPart
	err        error
}

// MultiUpload 分块上传
func (c *client) MultiUpload(bucketName, objectName, contentType string, fileSize int64, reader types.Reader) (string, error) {
	cur := time.Now()
	ossBucket, err := c.client.Bucket(bucketName)
	if err != nil {
		fmt.Printf("oss.MultiUpload bucket err: %s\n", err)
		return "", err
	}
	bucketIns := bucket{bucket: ossBucket}
	imur, err := bucketIns.multiUploadInit(objectName, contentType)
	if err != nil {
		fmt.Printf("oss.MultiUpload multiUploadInit err: %s\n", err)
		return "", err
	}
	chunks, err := splitByPartSize(fileSize, c.partSize, c.maxPartNum)
	var wg sync.WaitGroup
	partUploadResCh := make(chan partRes, len(chunks))
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk oss.FileChunk) {
			defer wg.Done()
			uploadPartRes, err := bucketIns.uploadPart(imur, reader, chunk)
			if err != nil {
				fmt.Printf("bucketIns.uploadPart chunk:%+v, err:%s\n", chunk, err)
			}
			partUploadResCh <- partRes{
				uploadPart: uploadPartRes,
				err:        err,
			}
		}(chunk)
	}
	wg.Wait()
	close(partUploadResCh)
	var parts []oss.UploadPart
	for _partRes := range partUploadResCh {
		if _partRes.err != nil {
			// 取消块上传
			bucketIns.bucket.AbortMultipartUpload(imur)
			return "", err
		}
		parts = append(parts, _partRes.uploadPart)
	}
	// 步骤3：完成分片上传，指定文件读写权限为公共读。
	cmur, err := bucketIns.bucket.CompleteMultipartUpload(imur, parts)
	fmt.Printf("oss.MultiUpload cost: %dms\n", time.Now().Sub(cur)/time.Millisecond)
	if err != nil {
		fmt.Printf("oss.MultiUpload  CompleteMultipartUpload err: %s\n", err)
		return "", err
	}
	fmt.Printf("oss.MultiUpload res: %+v\n", cmur)

	return fmt.Sprintf(c.baseUrl, bucketName) + "/" + objectName, err
}

func (b *bucket) multiUploadInit(objectName, contentType string) (oss.InitiateMultipartUploadResult, error) {
	opts := []oss.Option{
		oss.ContentType(contentType),
		oss.ObjectACL(oss.ACLPublicRead),
		oss.StorageClass(oss.StorageIA),
	}
	imur, err := b.bucket.InitiateMultipartUpload(objectName, opts...)
	return imur, err
}

func (b *bucket) uploadPart(imur oss.InitiateMultipartUploadResult, reader io.ReaderAt, chunk oss.FileChunk) (oss.UploadPart, error) {
	cur := time.Now()
	fd := io.NewSectionReader(reader, chunk.Offset, chunk.Size)
	part, err := b.bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
	fmt.Printf("[%s] oss.uploadPart chunk:%+v cost: %dms\n", time.Now(), chunk, time.Now().Sub(cur)/time.Millisecond)

	return part, err
}

func splitByPartSize(fileSize, chunkSize, maxChunkNum int64) ([]oss.FileChunk, error) {
	var chunkNum = fileSize / chunkSize
	if chunkNum >= maxChunkNum {
		return nil, errors.New("Too many parts, please increase part size")
	}
	var chunks []oss.FileChunk
	var chunk oss.FileChunk
	for i := int64(0); i < chunkNum; i++ {
		chunk.Number = int(i + 1)
		chunk.Offset = i * chunkSize
		chunk.Size = chunkSize
		chunks = append(chunks, chunk)
	}
	if fileSize%chunkSize > 0 {
		chunk.Number = len(chunks) + 1
		chunk.Offset = int64(len(chunks)) * chunkSize
		chunk.Size = fileSize % chunkSize
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}
