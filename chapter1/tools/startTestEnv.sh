

# 运行程序
LISTEN_ADDRESS=127.0.0.1:12345 STORAGE_ROOT=C:/tmp go run server.go

# 测试
curl -v 10.29.2.1:12345/objects/test1 -XPUT -d "this is object test1"
curl -v 10.29.2.1:12345/objects/test1