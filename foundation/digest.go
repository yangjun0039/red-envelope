package foundation

import (
	"crypto/md5"
	"fmt"
	"crypto/sha1"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
)

// MD5For16Length用于产生16位MD5消息摘要值，字母小写
func MD5For16Length(raw string) string {
	return MD5(raw)[8:24]
}

// 用于产生32位MD5消息摘要值，字母小写
func MD5(raw string) string {
	rawData := []byte(raw)
	digestText := fmt.Sprintf("%x", md5.Sum(rawData))
	return digestText
}

// 用于产生SHA-1消息摘要值，字母小写
func SHA1(s string) string {
	//产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(s))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要。
	bs := h.Sum(nil)
	//SHA1 值经常以 16 进制输出，例如在 git commit 中。使用%x 来将散列结果格式化为 16 进制字符串。
	str := fmt.Sprintf("%x", bs)
	return str
}

// SHA256摘要，输出字符串经过标准BASE64编码
func SHA256(s ...string) string {
	var buffer bytes.Buffer
	for _, v := range s {
		buffer.WriteString(v)
	}
	h := sha256.New()
	h.Write([]byte(buffer.String()))
	bs := h.Sum([]byte("ATl"))
	return base64.StdEncoding.EncodeToString(bs)
}
