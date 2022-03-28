

# 初始化文件目录
for i in `seq 1 6`
do
    mkdir -p C:/tmp/$i/objects
    mkdir -p C:/tmp/$i/temp
    mkdir -p C:/tmp/$i/garbage
done

# 启动程序
export RABBITMQ_SERVER=amqp://test:test@localhost:5672
export ES_SERVER=localhost:9200

LISTEN_ADDRESS=10.29.1.1:12345 STORAGE_ROOT=C:/tmp/1 go run dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.2:12345 STORAGE_ROOT=C:/tmp/2 go run dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.3:12345 STORAGE_ROOT=C:/tmp/3 go run dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.4:12345 STORAGE_ROOT=C:/tmp/4 go run dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.5:12345 STORAGE_ROOT=C:/tmp/5 go run dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.6:12345 STORAGE_ROOT=C:/tmp/6 go run dataServer/dataServer.go &

LISTEN_ADDRESS=10.29.2.1:12345 go run apiServer/apiServer.go &
LISTEN_ADDRESS=10.29.2.2:12345 go run apiServer/apiServer.go &


# ./tools/startTestEnv.sh chapter5



