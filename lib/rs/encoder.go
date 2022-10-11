package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

type encoder struct {
	writers []io.Writer
	cache   []byte
	enc     reedsolomon.Encoder //用来做输入数据缓存的字节数组
}

func NewEncoder(writers []io.Writer) *encoder {
	//生成具有4个数据片和2个校验片的RS编码器enc
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &encoder{
		writers: writers,
		cache:   nil,
		enc:     enc,
	}
}

//RSPUTStream本身没有实现Write方法，因为内嵌了encoder,所有会调用该方法
//如果缓存中的数据不满32000字节就暂不刷新，等待Write方法写一次被调用
func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	current := 0
	//将p中待写入的数据以块的形式放入缓存
	for length != 0 {
		next := BLOCK_SIZE - len(e.cache)
		if next > length {
			next = length
		}
		e.cache = append(e.cache, p[current:current+next]...)
		//如果缓存已满就调用flush方法将缓存实际写入到writers中
		if len(e.cache) == BLOCK_SIZE {
			e.Flush()
		}
		current += next
		length -= next
	}
	return len(p), nil
}
func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	//将缓存中的数据切成4个数据片
	shards, _ := e.enc.Split(e.cache)
	//调用Encode生成两个校验片
	e.enc.Encode(shards)
	//将6个分片的数据一次写入Writers并清空缓存
	for i := range shards {
		e.writers[i].Write(shards[i])
	}
	e.cache = []byte{}
}
