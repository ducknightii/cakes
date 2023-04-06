package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

func main() {
	str := "ewoJImNyZWF0ZWRfdG1wIiA6IDE2NjIzNjg2MTcyNDksCgkiZGF0YSIgOiAKCXsKCQkiZGV2aWNlX25hbWUiIDogIr/Y1sbG9zEgLSC2wb+oxvcyMiIsCgkJImhfaWQiIDogMTMsCgkJInN0YXR1cyIgOiAib2ZmbGluZSIKCX0sCgkiZXZlbnRfaWQiIDogIjE2NjIzNjg2MTcyNDkiLAoJImV2ZW50X3R5cGUiIDogImRldmljZV9zdGF0dXMiLAoJInZlcnNpb24iIDogIjAuMS4wIgp9"
	//str := "ewoJImNyZWF0ZWRfdG1wIiA6IDE2NjIzNjc4Nzk3NTgsCgkiZGF0YSIgOiAKCXsKCQkiZGV2aWNlX25hbWUiIDogIlx1MDdkOFx1MDU4Nlx1MDFiNzEgLSBcdTA1ODFcdTA3ZThcdTAxYjcyMiIsCgkJImhfaWQiIDogMTMsCgkJInN0YXR1cyIgOiAib2ZmbGluZSIKCX0sCgkiZXZlbnRfaWQiIDogIjE2NjIzNjc4Nzk3NTgiLAoJImV2ZW50X3R5cGUiIDogImRldmljZV9zdGF0dXMiLAoJInZlcnNpb24iIDogIjAuMS4wIgp9Cg=="
	dataBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
	}
	reader := transform.NewReader(bytes.NewReader(dataBytes), simplifiedchinese.GBK.NewDecoder())

	d, e := ioutil.ReadAll(reader)
	if e != nil {
		panic(e)
	}
	var data = make(map[string]interface{})
	err = json.Unmarshal(d, &data)
	fmt.Printf("str:[%s] data:[%+v] err:%v", d, data, err)

	card := 1563314106
	fmt.Printf("\ncard:%d %b %b\n", card%(1<<30), card, card%(1<<30))

	//dataBytes = dataBytes[67:74]
	/*for k, b := range dataBytes {
		fmt.Printf("[%d], byte:[%b][%d] s:%s\n", k, b, int(b), string(b))
	}

	c := dataBytes[0:2]
	bytebuff := bytes.NewBuffer([]byte{
		111,
	})
	var data, data2 = new(int64), new(int64)
	binary.Read(bytebuff, binary.LittleEndian, &data2)
	binary.Read(bytebuff, binary.BigEndian, &data)

	fmt.Printf("[%x][%d][%d][%s]", c, data, data2, string(c))*/
}
