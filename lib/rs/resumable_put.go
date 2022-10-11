package rs

import (
	"OSS/lib/objectStream"
	"OSS/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type resumableToken struct {
	Name    string   //对象名称
	Size    int64    //对象内容的大小
	Hash    string   //对象内容的hash
	Servers []string //6个分片所在的地址
	Uuids   []string //6个分片的uuid
}
type RSResumablePutStream struct {
	*RSPutStream
	*resumableToken
}

func NewRSResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, err := NewRSPutStream(dataServers, hash, size)
	if err != nil {
		return nil, err
	}
	uuids := make([]string, ALL_SHARDS)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectStream.TempPutStream).Uuid
	}
	token := &resumableToken{
		Name:    name,
		Size:    size,
		Hash:    hash,
		Servers: dataServers,
		Uuids:   uuids,
	}
	return &RSResumablePutStream{
		RSPutStream:    putStream,
		resumableToken: token,
	}, nil
}

//将自身数据以json格式编码，经过base64编码后返回
func (rsr *RSResumablePutStream) ToToken() string {
	bytes, _ := json.Marshal(rsr.resumableToken)
	return base64.StdEncoding.EncodeToString(bytes)
}
func NewRSResumablePutStreamFromToken(ftoken string) (*RSResumablePutStream, error) {
	token, err := url.PathUnescape(ftoken)
	if err != nil {
		return nil, err
	}
	//对token进行base64的解码
	bytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	var t resumableToken
	err = json.Unmarshal(bytes, &t)
	if err != nil {
		return nil, err
	}
	writers := make([]io.Writer, ALL_SHARDS)
	for i := range writers {
		writers[i] = &objectStream.TempPutStream{Server: t.Servers[i], Uuid: t.Uuids[i]}
	}
	enc := NewEncoder(writers)
	return &RSResumablePutStream{
		RSPutStream:    &RSPutStream{enc},
		resumableToken: &t,
	}, nil
}

//获取数据节点已经储存该对象多少数据了
func (rsr *RSResumablePutStream) CurrentSize() int64 {
	//以head方法获取第一个分片临时对象的大小
	res, err := http.Head(fmt.Sprintf("http://%s/temp/%s", rsr.Servers[0], rsr.Uuids[0]))
	if err != nil {
		log.Println(err)
		return -1
	}
	if res.StatusCode != http.StatusOK {
		log.Println(res.StatusCode)
		return -1
	}
	size := utils.GetSizeFromHeader(res.Header) * DATA_SHARDS
	if size > rsr.Size {
		size = rsr.Size
	}
	return size
}
