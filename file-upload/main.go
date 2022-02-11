package main

// 参考 https://tonybai.com/2021/01/16/upload-and-download-file-using-multipart-form-over-http/

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/ducknightii/cakes/file-upload/storage"
)

func main() {
	storage.Init()
	http.HandleFunc("/upload", handleUploadFile)
	fmt.Println("listen: 18800")
	http.ListenAndServe(":18800", nil)
}

func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	memStat()
	r.ParseMultipartForm(1 * 1024 * 1024) // 1Mi
	mForm := r.MultipartForm
	if mForm == nil {
		return
	}
	// todo cancel deal
	kvs := mForm.Value
	var method string = "bytes"
	if len(kvs) > 0 && len(kvs["method"]) > 0 {
		method = kvs["method"][0]
	}
	for k, _ := range mForm.File {
		// k is the key of file part
		file, fileHeader, err := r.FormFile(k)
		if err != nil {
			fmt.Println("inovke FormFile error:", err)
			return
		}
		defer file.Close()
		fmt.Printf("the uploaded file: name[%s], size[%d], header[%#v]\n",
			fileHeader.Filename, fileHeader.Size, fileHeader.Header)

		memStat()

		// store uploaded file into local path
		localFile := fileHeader.Filename
		if method == "bytes" {
			err = fileSaveBytes(file, localFile)
		} else if method == "multi" {
			err = fileSaveMulti(file, localFile, fileHeader.Size)
		} else {
			err = fileSaveStream(file, localFile)
		}
		if err != nil {
			fmt.Printf("file %s uploaded err: %s\n", fileHeader.Filename, err)
			continue
		}

		fmt.Printf("file %s uploaded ok\n", fileHeader.Filename)

		memStat()
	}

	// notes: 因为计算内存 所以没用defer 正常还是defer吧
	r.MultipartForm.RemoveAll()
	runtime.GC()
	time.Sleep(time.Millisecond * 200)
	fmt.Println("deal finish")
	memStat()
}

func memStat() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Mem Alloc: %dK Free: %dK  HeapInuse: %dK StackInuse: %dK\n", m.Alloc/1024, m.Frees/1024, m.HeapInuse/1024, m.StackInuse/1024)
}

func summary(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

func fileSaveBytes(file multipart.File, localFile string) error {
	memStat()
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	fmt.Printf("file copy finish, len:%dK cap: %dK\n", len(buf.Bytes())/1024, cap(buf.Bytes())/1024)
	memStat()
	if err != nil {
		return err
	}

	/*// 一个测试 newbuffer 不会直接分配内存
	_ = bytes.NewBuffer(buf.Bytes())
	fmt.Println("bytes.NewBufferfinish")
	memStat()
	// 只增加 io.copy 的32K
	_, _ = bce.NewBodyFromBytes(buf.Bytes())
	fmt.Println("bce.NewBodyFromBytes")
	memStat()*/

	saveUrl, err := storage.Storage.Upload("***", localFile, "application/octet-stream", buf.Bytes())
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 3) // bos 传完
	fmt.Printf("SaveUrl: %s\n", saveUrl)
	memStat()
	// 光标跳到开头
	file.Seek(0, 0)
	hash := md5.New()
	io.Copy(hash, file)
	fmt.Println("hash io copy finish")

	memStat()

	if err == nil {
		fmt.Println("Hex:", hex.EncodeToString(hash.Sum(nil)))
		memStat()
	}
	return err
}

func fileSaveStream(file multipart.File, localFile string) error {
	saveUrl, err := storage.Storage.UploadFromStream("***", localFile, "application/octet-stream", file)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 3) // bos 传完
	fmt.Printf("SaveUrl: %s\n", saveUrl)
	memStat()
	// 光标跳到开头
	file.Seek(0, 0)
	hash := md5.New()
	io.Copy(hash, file)
	fmt.Println("hash io copy finish")

	memStat()

	if err == nil {
		fmt.Println("Hex:", hex.EncodeToString(hash.Sum(nil)))
		memStat()
	}
	return err
}

func fileSaveMulti(file multipart.File, localFile string, fileSize int64) error {
	saveUrl, err := storage.Storage.MultiUpload("***", localFile, "application/octet-stream", fileSize, file)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 3) // bos 传完
	fmt.Printf("SaveUrl: %s\n", saveUrl)
	memStat()
	// 光标跳到开头
	file.Seek(0, 0)
	hash := md5.New()
	io.Copy(hash, file)
	fmt.Println("hash io copy finish")

	memStat()

	if err == nil {
		fmt.Println("Hex:", hex.EncodeToString(hash.Sum(nil)))
		memStat()
	}
	return err
}

// 边读边写 边计算hash
// todo 融合进storage：从src fd中读取数据 同时完成两个独立功能 storage存储 和 hash计算，而不是在一个方法里处理, 是否可以用io.Pipe 构造两个ch 一个用于存储一个用于hash计算
func fileSave2(file multipart.File, localFile string) error {
	out, err := os.Create(localFile)
	if err != nil {
		fmt.Printf("failed to open the file %s for writing", localFile)
		return err
	}
	defer out.Close()

	hash := md5.New()
	buf := make([]byte, 32*1024)
	for {
		memStat()
		nr, er := file.Read(buf)
		if nr > 0 {
			// todo err handle
			hash.Write(buf[0:nr])
			nw, ew := out.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write result")
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	fmt.Println("io finish")
	memStat()
	if err == nil {
		fmt.Println("Hex:", hex.EncodeToString(hash.Sum(nil)))
		memStat()
	}

	return err
}
