# 模拟6 dataserver 2 apiserver

sudo ifconfig lo:1 10.29.1.1/16
sudo ifconfig lo:2 10.29.1.2/16
sudo ifconfig lo:3 10.29.1.3/16
sudo ifconfig lo:4 10.29.1.4/16
sudo ifconfig lo:5 10.29.1.5/16
sudo ifconfig lo:6 10.29.1.6/16
sudo ifconfig lo:8 10.29.2.1/16
sudo ifconfig lo:9 10.29.2.2/16


# 分配文件夹
for i in `seq 1 6`
do
    rm -rf /tmp/$i/objects/*
    rm -rf /tmp/$i/temp/*
    rm -rf /tmp/$i/garbage/*
done

for i in `seq 1 6`
do
    mkdir -p /log/$i/
done

mkdir -p /log/21
mkdir -p /log/22

# 环境变量

export ES_SERVER=127.0.0.1:9200
export REDIS_CLUSTER=127.0.0.1:6379
export REDIS_PASSWORD=

# 运行
killall apiServer
killall dataServer

export ES_SERVER=127.0.0.1:9200
LOG_DIRECTORY=/log/1 LISTEN_ADDRESS=10.29.1.1:12345 STORAGE_ROOT=/tmp/1 go run dataServer/dataServer.go &
LOG_DIRECTORY=/log/2 LISTEN_ADDRESS=10.29.1.2:12345 STORAGE_ROOT=/tmp/2 go run dataServer/dataServer.go &
LOG_DIRECTORY=/log/3 LISTEN_ADDRESS=10.29.1.3:12345 STORAGE_ROOT=/tmp/3 go run dataServer/dataServer.go &
LOG_DIRECTORY=/log/4 LISTEN_ADDRESS=10.29.1.4:12345 STORAGE_ROOT=/tmp/4 go run dataServer/dataServer.go &
LOG_DIRECTORY=/log/5 LISTEN_ADDRESS=10.29.1.5:12345 STORAGE_ROOT=/tmp/5 go run dataServer/dataServer.go &
LOG_DIRECTORY=/log/6 LISTEN_ADDRESS=10.29.1.6:12345 STORAGE_ROOT=/tmp/6 go run dataServer/dataServer.go &
LOG_DIRECTORY=/log/21 LISTEN_ADDRESS=10.29.2.1:12345 go run apiServer/apiServer.go &
LOG_DIRECTORY=/log/22 LISTEN_ADDRESS=10.29.2.2:12345 go run apiServer/apiServer.go &
