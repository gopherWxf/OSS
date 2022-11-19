package rs

/*
	采用4+2的RS码策略，存储空间150%，可以丢失6个分片对象中的任意两个
	M+N的RS码
 	储存空间：(M+N)/M*100%
	抵御能力：N（可以丢失的分片数量）
*/
const (
	DATA_SHARDS     = 4                             //数据片 M
	PARITY_SHARDS   = 2                             //校验片 N
	ALL_SHARDS      = DATA_SHARDS + PARITY_SHARDS   //所有分片
	BLOCK_PER_SHARD = 1310720                       //块碎片大小    1.25M*4=5M
	BLOCK_SIZE      = BLOCK_PER_SHARD * DATA_SHARDS //块大小=数据片*块碎片大小
)
