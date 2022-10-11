package temp

import (
	"OSS/apiServer/locate"
	es "OSS/lib/ElasticSearch"
	"OSS/lib/rs"
	"OSS/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//获取token
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	//通过token获得RSResumablePutStream的结构体指针
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//获取数据节点已经储存该对象多少数据了
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//从range头部获得offset
	offset := utils.GetOffsetFromHeader(r.Header)
	//如果不一致则返回错误
	if current != offset {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	//如果一致，则在for循环中以32000字节长度读取正文并写入stream,分块
	bytes := make([]byte, rs.BLOCK_SIZE)
	for {
		n, err := io.ReadFull(r.Body, bytes)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current += int64(n)
		//如果读到的总长度超过对象的大小，说明客户端上传额数据有误
		if current > stream.Size {
			//删除临时对象
			stream.Commit(false)
			log.Println("resumable put exceed size")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		//n != rs.BLOCK_SIZE：说明本次客户端的数据已经传完了
		//current != stream.Size：说明对象的整体数据还没有完全传输完
		//此时接口服务会丢弃这次读取的长度不到32000字节的数据
		if n != rs.BLOCK_SIZE && current != stream.Size {
			return
		}
		stream.Write(bytes[:n])
		//如果读到的总长度等于对象的大小，说明客户端上传了全部数据
		if current == stream.Size {
			//调用flush方法将剩余数据写进临时对象
			stream.Flush()
			//调用rs.NewRSResumableGetStream生成一个临时对象读取流
			getStream, _ := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			//读取流中的数据并计算hash值
			hash := url.PathEscape(utils.CalculateHash(getStream))
			//如果hash值不一致，则说明数据有误，删除临时对象
			if hash != stream.Hash {
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			//如果hash一致，检查是否已经存在，存在则删除，不存在则转正
			if locate.Exist(hash) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			realhash, err := url.PathUnescape(stream.Hash)
			if err != nil {
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			//添加进元数据es
			err = es.AddVersion(stream.Name, realhash, stream.Size)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}
