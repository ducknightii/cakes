package bos

// https://cloud.baidu.com/doc/BOS/s/Djwvyrxn6#%E4%B8%8A%E4%BC%A0%E6%96%87%E4%BB%B6

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
)

var (
	once sync.Once
	ins  = new(client)
)

type client struct {
	client  *bos.Client
	baseUrl string
}

// todo 分块上传功能 完成
type multiUploadIns struct {
	SaveUrl string

	client     *client
	bucketName string
	objectName string
	uploadID   string
	partEtags  etagSlice
}
type etagSlice []api.UploadInfoType

func (p etagSlice) Len() int           { return len(p) }
func (p etagSlice) Less(i, j int) bool { return p[i].PartNumber < p[j].PartNumber }
func (p etagSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type partRes struct {
	etag api.UploadInfoType
	err  error
}

func Init(ak, sk, endpoint, baseUrl string) (c *client, err error) {
	once.Do(func() {
		ins.client, err = bos.NewClient(ak, sk, endpoint)
		ins.baseUrl = baseUrl
		ins.client.MultipartSize = 20 * bos.MULTIPART_ALIGN // 20Mi
		ins.client.MaxParallel = 100

	})
	return ins, err
}

func (c *client) Upload(bucketName, objectName, contentType string, body []byte) (string, error) {
	cur := time.Now()
	args := new(api.PutObjectArgs)
	args.StorageClass = api.STORAGE_CLASS_STANDARD_IA
	args.ContentType = contentType
	args.ContentLength = int64(len(body))
	_, err := c.client.PutObjectFromBytes(bucketName, objectName, body, args)
	fmt.Printf("bos.Upload cost: %dms\n", time.Now().Sub(cur)/time.Millisecond)
	if err != nil {
		fmt.Println("bos-err", fmt.Sprintf("bucket:%s object:%s err:%s", bucketName, objectName, err))
	}

	return fmt.Sprintf(c.baseUrl, bucketName) + "/" + objectName, err
}

// UploadFromStream 流式上传
// notes: 任然是全部读取到内存再发起请求，由于用的io.reader 关闭由外部控制，所以采用阻塞上传，避免还没读完被关闭
// 内存的耗费 io.Reader -> buffer -> body 大概翻了4倍 (参考slice扩容规律)  (还未完全读完具体实现)
func (c *client) UploadFromStream(bucketName, objectName, contentType string, reader io.Reader) (string, error) {
	cur := time.Now()
	args := new(api.PutObjectArgs)
	args.StorageClass = api.STORAGE_CLASS_STANDARD_IA
	args.ContentType = contentType
	_, err := c.client.PutObjectFromStream(bucketName, objectName, reader, args)
	fmt.Printf("bos.UploadFromStream cost: %dms\n", time.Now().Sub(cur)/time.Millisecond)
	if err != nil {
		fmt.Println("bos-err", fmt.Sprintf("bucket:%s object:%s err:%s", bucketName, objectName, err))
	}

	return fmt.Sprintf(c.baseUrl, bucketName) + "/" + objectName, err
}

// MultiUpload 分块上传
func (c *client) MultiUpload(bucketName, objectName, contentType string, fileSize int64, reader io.Reader) (string, error) {
	cur := time.Now()

	partSize := (c.client.MultipartSize +
		bos.MULTIPART_ALIGN - 1) / bos.MULTIPART_ALIGN * bos.MULTIPART_ALIGN

	if fileSize <= partSize {
		return "", errors.New("no need multi upload")
	}
	multiIns, err := c.multiUploadInit(bucketName, objectName, contentType)
	if err != nil {
		return "", err
	}
	partNum := (fileSize + partSize - 1) / partSize
	if partNum > bos.MAX_PART_NUMBER {
		return "", errors.New("file too large")
	}

	var wg sync.WaitGroup
	partUploadResCh := make(chan partRes, partNum)
	for i := int64(1); i <= partNum; i++ {
		// 计算偏移offset和本次上传的大小uploadSize
		uploadSize := partSize
		offset := partSize * (i - 1)
		left := fileSize - offset
		if left < partSize {
			uploadSize = left
		}

		buf := new(bytes.Buffer)
		_, err = io.CopyN(buf, reader, uploadSize)
		if err != nil {
			fmt.Printf("partNo:%d io.Copy err:%s\n", i, err)
			partUploadResCh <- partRes{
				etag: api.UploadInfoType{},
				err:  fmt.Errorf("partNo:%d io copy err %w", i, err),
			}
			continue
		}

		wg.Add(1)
		go func(partNo int, body []byte) {
			defer wg.Done()
			etag, err := multiIns.uploadPart(partNo, body)
			fmt.Printf("partNo:%d upload err:%v\n", partNo, err)
			partUploadResCh <- partRes{
				etag: etag,
				err:  err,
			}
		}(int(i), buf.Bytes())
	}
	wg.Wait()
	close(partUploadResCh)
	for _partRes := range partUploadResCh {
		if _partRes.err != nil {
			// 取消块上传
			c.client.AbortMultipartUpload(multiIns.bucketName, multiIns.objectName, multiIns.uploadID)
			return "", err
		}
		multiIns.partEtags = append(multiIns.partEtags, _partRes.etag)
	}
	sort.Sort(multiIns.partEtags)
	fmt.Printf("multiIns.partEtags r%#v\n", multiIns.partEtags)
	res, err := multiIns.complete()
	fmt.Printf("complete res: %#v err:%v\n", res, err)
	fmt.Printf("bos.MultiUpload cost: %dms\n", time.Now().Sub(cur)/time.Millisecond)

	return fmt.Sprintf(c.baseUrl, bucketName) + "/" + objectName, err
}

// MultiUploadInit 分块上传初始化一个序号uploadID  ep: contentType=application/octet-stream
func (c *client) multiUploadInit(bucketName, objectName, contentType string) (ins multiUploadIns, err error) {
	/*res, err := c.client.BasicInitiateMultipartUpload(bucketName, objectName)
	if err != nil {
		fmt.Println("BasicInitiateMultipartUpload err ", err)
		return
	}
	fmt.Printf("BasicInitiateMultipartUpload: %#v\n", res)*/

	args := new(api.InitiateMultipartUploadArgs)
	args.StorageClass = api.STORAGE_CLASS_STANDARD_IA
	res, err := c.client.InitiateMultipartUpload(bucketName, objectName, contentType, args)
	if err != nil {
		fmt.Println("InitiateMultipartUpload err ", err)
		return
	}
	fmt.Printf("InitiateMultipartUpload: %#v\n", res)
	return multiUploadIns{
		SaveUrl:    fmt.Sprintf(c.baseUrl, bucketName) + "/" + objectName,
		client:     c,
		bucketName: bucketName,
		objectName: objectName,
		uploadID:   res.UploadId,
		partEtags:  make([]api.UploadInfoType, 0),
	}, nil
}

// MultiUploadPart 上传一块数据
func (multiIns *multiUploadIns) uploadPart(partNo int, body []byte) (api.UploadInfoType, error) {
	cur := time.Now()
	// 上传当前分块
	etag, err := multiIns.client.client.UploadPartFromBytes(multiIns.bucketName, multiIns.objectName, multiIns.uploadID, partNo, body, nil)

	fmt.Printf("[%s] bos.uploadPart len:%d cost: %dms\n", time.Now(), len(body), time.Now().Sub(cur)/time.Millisecond)

	if err != nil {
		fmt.Printf("uploadPart partNo:%d err:%s\n", partNo, err)
		return api.UploadInfoType{}, err
	}

	return api.UploadInfoType{
		PartNumber: partNo,
		ETag:       etag,
	}, nil
}

func (multiIns *multiUploadIns) complete() (*api.CompleteMultipartUploadResult, error) {
	completeArgs := api.CompleteMultipartUploadArgs{Parts: multiIns.partEtags}
	return multiIns.client.client.CompleteMultipartUploadFromStruct(
		multiIns.bucketName, multiIns.objectName, multiIns.uploadID, &completeArgs)
}
