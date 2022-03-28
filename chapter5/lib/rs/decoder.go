package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

type decoder struct {
	readers   []io.Reader
	writers   []io.Writer
	enc       reedsolomon.Encoder
	size      int64  //对象的大小
	cache     []byte //缓存读取的数据
	cacheSize int
	total     int64 //表示已读多少字节
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64) *decoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &decoder{
		readers:   readers,
		writers:   writers,
		enc:       enc,
		size:      size,
		cache:     nil,
		cacheSize: 0,
		total:     0,
	}
}
func (d *decoder) Read(p []byte) (n int, err error) {
	//当cache中没有数据的时候，会调用getData获取数据
	if d.cacheSize == 0 {
		err = d.getData()
		if err != nil {
			return 0, err
		}
	}
	//如果length超过缓存数据大小，则调整
	length := len(p)
	if d.cacheSize < length {
		length = d.cacheSize
	}
	d.cacheSize -= length
	copy(p, d.cache[:length])
	//调整缓存，仅保留剩下的部分
	d.cache = d.cache[length:]
	return length, nil
}
func (d *decoder) getData() error {
	//如果已经读的数据与对象大小一致，说明所有数据都读完了
	if d.total == d.size {
		return io.EOF
	}
	//6个分片
	shards := make([][]byte, ALL_SHARDS)
	//修复分片
	repairIds := make([]int, 0)
	for i := range shards {
		//如果是nil，说明数据丢失，存入待恢复数组
		if d.readers[i] == nil {
			repairIds = append(repairIds, i)
		} else {
			//读取流正常，读取数据
			shards[i] = make([]byte, BLOCK_PER_SHARD)
			n, err := io.ReadFull(d.readers[i], shards[i])
			//如果发生非EOF失败，则将shards置为nil
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				shards[i] = nil
			} else if n != BLOCK_PER_SHARD {
				shards[i] = shards[i][:n]
			}
		}
	}
	//尝试将nil的shards恢复
	err := d.enc.Reconstruct(shards)
	//如果这一步错误，说明对象不可被修复了
	if err != nil {
		return err
	}
	//将数据写入将需要恢复的分片的writer
	for i := range repairIds {
		id := repairIds[i]
		d.writers[id].Write(shards[id])
	}
	//最后，遍历4个数据分片，将数据添加到缓存
	for i := 0; i < DATA_SHARDS; i++ {
		shardSize := int64(len(shards[i]))
		if d.total+shardSize > d.size {
			shardSize -= d.total + shardSize - d.size
		}
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.cacheSize += int(shardSize)
		d.total += shardSize
	}
	return nil
}
