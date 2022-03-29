package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//获取请求头部的hash值
func GetHashFromHeader(h http.Header) string {
	digest := h.Get("digest")
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}

//获取内容的长度
func GetSizeFromHeader(h http.Header) int64 {
	size, _ := strconv.ParseInt(h.Get("content-length"), 0, 64)
	return size
}

//计算内容的hash值
func CalculateHash(r io.Reader) string {
	//生成sha256.digest结构体的实例，其实现拉hash.Hash接口
	h := sha256.New()
	//将r写入h，h会对写入的数据进行计算
	io.Copy(h, r)
	//h.Sum返回hash的二进制数据，用base64进行编码
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
func GetOffsetFromHeader(h http.Header) int64 {
	byteRange := h.Get("range")
	if len(byteRange) < 7 {
		return 0
	}
	if byteRange[:6] != "bytes=" {
		return 0
	}
	bytePos := strings.Split(byteRange[6:], "-")
	offset, _ := strconv.ParseInt(bytePos[0], 0, 64)
	return offset
}
