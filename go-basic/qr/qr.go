package main

import (
	"bytes"
	"crypto/des"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func main() {
	//encStr := GenEncHexCode(11, int(time.Now().Unix()), int(time.Now().Unix()))
	//fmt.Println("encStr:", encStr)
	//
	//encBytes, _ := hex.DecodeString(encStr)
	decData := desECBDecrypt([]byte("DC186D5D39F6EC1ED57DD816937F08A6CB7CD966D003F45068A0D8FF6E36119D68A0D8FF6E36119D269195CEF1C491CE"))
	fmt.Println("decData: ", decData)
	fmt.Println("decData: ", decRawData(decData))
}

func GenEncHexCode(card, startTs, endTs int) string {
	data := genRawData(card, startTs, endTs)

	encRes := desECBEncrypt(data)

	encStr := strings.ToUpper(hex.EncodeToString(encRes))

	return encStr
}

func genRawData(card, startTs, endTs int) []byte {
	now := int(time.Now().Unix())

	buffer := new(bytes.Buffer)
	buffer.WriteByte(0xD0)
	binary.Write(buffer, binary.BigEndian, int32(now))
	buffer.WriteByte(0x28) //数据包长度

	buffer.Write([]byte{0x00, 0x01})                       //项目ID
	binary.Write(buffer, binary.BigEndian, int32(card))    //卡号
	buffer.WriteByte(0xA1)                                 //权限类型
	binary.Write(buffer, binary.BigEndian, int32(startTs)) //生效时间
	binary.Write(buffer, binary.BigEndian, int32(endTs))   //无效时间
	for i := 0; i < 24; i++ {
		buffer.WriteByte(0x00)
	}
	data := buffer.Bytes()
	fmt.Println("data:", data)
	sign := data[6]
	for i := 7; i < 45; i++ {
		sign = sign ^ data[i]
	}
	buffer.WriteByte(sign)
	data = buffer.Bytes()
	md5Row := strings.ToUpper(hex.EncodeToString(data[6:46])) + "@" + strings.ToUpper(hex.EncodeToString(data[1:5]))
	md := md5.Sum([]byte(md5Row))
	buffer.Write(md[len(md)-2:])

	out := buffer.Bytes()
	out[0] = 0xD0 + sign&0x0F
	fmt.Printf("card:%d [%d-%d] gen res:[%s]\n", card, startTs, endTs, strings.ToUpper(hex.EncodeToString(out)))

	return out
}

func decRawData(data []byte) int32 {
	cardBytes := data[8:12]

	fmt.Println("cardBytes:", cardBytes)

	bytebuff := bytes.NewBuffer(cardBytes)

	var card = new(int32)

	fmt.Println(bytebuff.Len())

	binary.Read(bytebuff, binary.BigEndian, card)

	return *card

}

func desECBEncrypt(data []byte) []byte {
	key, _ := hex.DecodeString("B6FECEACC3DCD4BF")

	block, err := des.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bs := block.BlockSize()

	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out
}

func desECBDecrypt(data []byte) []byte {
	key, _ := hex.DecodeString("B6FECEACC3DCD4BF")

	block, err := des.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bs := block.BlockSize()

	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out
}
