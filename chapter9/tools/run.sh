
# 上传一个小对象
curl -v 10.29.2.1:12345/objects/test3 -XPUT -d "this is object test3 version 1" -H "Digest: SHA-256=p3NoRweZQOBNCrQa0QuCVKV2RiDWEst1/G+O8dxafq0="
# 上传同名对象，修改其内容
curl -v 10.29.2.1:12345/objects/test3 -XPUT -d "this is object test3 version 2" -H "Digest: SHA-256=cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo="
# 查看对象，发现是最新数据
curl -v 10.29.2.1:12345/objects/test3
# 查看对象是历史数据
curl -v 10.29.2.1:12345/objects/test3?version=1
# 查看对象保存在哪台数据节点（对于先版本已经没用了，因为做了冗余备份，数据存放在6台节点上）
curl -v 10.29.2.1:12345/locate/cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo=
# 删除对象
curl -v 10.29.2.1:12345/objects/test3 -XDELETE
# 上传对象
curl -v 10.29.2.1:12345/objects/test4 -XPUT -d "this is object test4" -H "Digest: SHA-256=qxn98QM9ZPosMkBIKOwQUTI+5s2a4sDNaBBlTT5jLhw="
# 上传对象
curl -v 10.29.2.1:12345/objects/test5 -XPUT -d "this is object test5" -H "Digest: SHA-256=B494C1vj+98Y+PTGRiNqWu7gRgWQwiHnEofa47sN6mk="
# 查看所有版本信息
curl -v http://10.29.2.1:12345/versions/

#-----------上传大对象(分批次上传)-----------#
# 这里演示一次性上传全部
# 1.首先使用POST接口，告诉服务器，客户端要上传的总的长度与hash值。服务器会返回给客户端一个token，待对象上传完成后，再去检查长度与hash值是否一致
curl 10.29.2.1:12345/objects/test6  -v -XPOST -H "Digest: SHA-256=kZLCW3NPy62+MtrcKAicYNsOOfkMwgzi5XM/VyYazAw=" -H "size: 100000"
# 2.使用HEAD接口，查询当前在这个token上面真实的被接收上传了多少字节
curl -v -XHEAD 10.29.2.1:12345/temp/eyJOYW1lIjoidGVzdDYiLCJTaXplIjoxMDAwMDAsIkhhc2giOiJrWkxDVzNOUHk2MitNdHJjS0FpY1lOc09PZmtNd2d6aTVYTSUyRlZ5WWF6QXc9IiwiU2VydmVycyI6WyIxMC4yOS4xLjI6MTIzNDUiLCIxMC4yOS4xLjQ6MTIzNDUiLCIxMC4yOS4xLjU6MTIzNDUiLCIxMC4yOS4xLjE6MTIzNDUiLCIxMC4yOS4xLjM6MTIzNDUiLCIxMC4yOS4xLjY6MTIzNDUiXSwiVXVpZHMiOlsiN2JkMGQ5ZWQtYWY2OC0xMWVjLTllZTEtYTIzNjk1OTM3MTQ0IiwiN2JkMTJhM2QtYWY2OC0xMWVjLWFlZDQtYTIzNjk1OTM3MTQ0IiwiN2JkMThlMTEtYWY2OC0xMWVjLTljY2ItYTIzNjk1OTM3MTQ0IiwiN2JkMWNhNDYtYWY2OC0xMWVjLThkYmYtYTIzNjk1OTM3MTQ0IiwiN2JkMjFhM2ItYWY2OC0xMWVjLTk2OTgtYTIzNjk1OTM3MTQ0IiwiN2JkMjZhMjctYWY2OC0xMWVjLThkYjctYTIzNjk1OTM3MTQ0Il19
# 2.使用PUT接口，上传某一部分，所以这个请求需要提供一个range（不写默认从0开始）
curl -v -XPUT  --data-binary @C:/Users/68725/Desktop/file  10.29.2.1:12345/temp/eyJOYW1lIjoidGVzdDYiLCJTaXplIjoxMDAwMDAsIkhhc2giOiJrWkxDVzNOUHk2MitNdHJjS0FpY1lOc09PZmtNd2d6aTVYTSUyRlZ5WWF6QXc9IiwiU2VydmVycyI6WyIxMC4yOS4xLjI6MTIzNDUiLCIxMC4yOS4xLjQ6MTIzNDUiLCIxMC4yOS4xLjU6MTIzNDUiLCIxMC4yOS4xLjE6MTIzNDUiLCIxMC4yOS4xLjM6MTIzNDUiLCIxMC4yOS4xLjY6MTIzNDUiXSwiVXVpZHMiOlsiN2JkMGQ5ZWQtYWY2OC0xMWVjLTllZTEtYTIzNjk1OTM3MTQ0IiwiN2JkMTJhM2QtYWY2OC0xMWVjLWFlZDQtYTIzNjk1OTM3MTQ0IiwiN2JkMThlMTEtYWY2OC0xMWVjLTljY2ItYTIzNjk1OTM3MTQ0IiwiN2JkMWNhNDYtYWY2OC0xMWVjLThkYmYtYTIzNjk1OTM3MTQ0IiwiN2JkMjFhM2ItYWY2OC0xMWVjLTk2OTgtYTIzNjk1OTM3MTQ0IiwiN2JkMjZhMjctYWY2OC0xMWVjLThkYjctYTIzNjk1OTM3MTQ0Il19
# 3.当文件上传完成之后，token失效，不能再用第二步来看发送多少字节了
# token
/temp/eyJOYW1lIjoidGVzdDYiLCJTaXplIjoxMDAwMDAsIkhhc2giOiJrWkxDVzNOUHk2MitNdHJjS0FpY1lOc09PZmtNd2d6aTVYTSUyRlZ5WWF6QXc9IiwiU2VydmVycyI6WyIxMC4yOS4xLjI6MTIzNDUiLCIxMC4yOS4xLjQ6MTIzNDUiLCIxMC4yOS4xLjU6MTIzNDUiLCIxMC4yOS4xLjE6MTIzNDUiLCIxMC4yOS4xLjM6MTIzNDUiLCIxMC4yOS4xLjY6MTIzNDUiXSwiVXVpZHMiOlsiN2JkMGQ5ZWQtYWY2OC0xMWVjLTllZTEtYTIzNjk1OTM3MTQ0IiwiN2JkMTJhM2QtYWY2OC0xMWVjLWFlZDQtYTIzNjk1OTM3MTQ0IiwiN2JkMThlMTEtYWY2OC0xMWVjLTljY2ItYTIzNjk1OTM3MTQ0IiwiN2JkMWNhNDYtYWY2OC0xMWVjLThkYmYtYTIzNjk1OTM3MTQ0IiwiN2JkMjFhM2ItYWY2OC0xMWVjLTk2OTgtYTIzNjk1OTM3MTQ0IiwiN2JkMjZhMjctYWY2OC0xMWVjLThkYjctYTIzNjk1OTM3MTQ0Il19


# 这里演示分两次上传
# 1.首先使用POST接口，告诉服务器，客户端要上传的总的长度与hash值。服务器会返回给客户端一个token，待对象上传完成后，再去检查长度与hash值是否一致
curl 10.29.2.1:12345/objects/test6  -v -XPOST -H "Digest: SHA-256=kZLCW3NPy62+MtrcKAicYNsOOfkMwgzi5XM/VyYazAw=" -H "size: 100000"
# 2.first 有50%的数据，但是不是32000的倍数，也没有传输全部数据，所以会有部分数据失效
curl -v -XPUT  --data-binary @C:/Users/68725/Desktop/first  10.29.2.1:12345
curl -v -XPUT  --data-binary @C:/Users/68725/Desktop/first  10.29.2.1:12345/temp/eyJOYW1lIjoidGVzdDYiLCJTaXplIjoxMDAwMDAsIkhhc2giOiJrWkxDVzNOUHk2MitNdHJjS0FpY1lOc09PZmtNd2d6aTVYTSUyRlZ5WWF6QXc9IiwiU2VydmVycyI6WyIxMC4yOS4xLjI6MTIzNDUiLCIxMC4yOS4xLjQ6MTIzNDUiLCIxMC4yOS4xLjM6MTIzNDUiLCIxMC4yOS4xLjY6MTIzNDUiLCIxMC4yOS4xLjE6MTIzNDUiLCIxMC4yOS4xLjU6MTIzNDUiXSwiVXVpZHMiOlsiYWZkYzFiYzItYWY2Zi0xMWVjLThkMjEtYTIzNjk1OTM3MTQ0IiwiYWZlMDMxNTYtYWY2Zi0xMWVjLWJjZDUtYTIzNjk1OTM3MTQ0IiwiYWZlNGJlMjEtYWY2Zi0xMWVjLWIxNTYtYTIzNjk1OTM3MTQ0IiwiYWZlOTk0MDUtYWY2Zi0xMWVjLTgyMDItYTIzNjk1OTM3MTQ0IiwiYWZlZTg0YzgtYWY2Zi0xMWVjLWFmNzMtYTIzNjk1OTM3MTQ0IiwiYWZmMmQ1MWYtYWY2Zi0xMWVjLTk1YzEtYTIzNjk1OTM3MTQ0Il19

# 使用HEAD接口查看上面真实上传了多少数据   Content-Length: 32000 实际32000，剩下的18000由于没有达到块大小，被服务器丢弃了
curl -I 10.29.2.1:12345
curl -I 10.29.2.1:12345/temp/eyJOYW1lIjoidGVzdDYiLCJTaXplIjoxMDAwMDAsIkhhc2giOiJrWkxDVzNOUHk2MitNdHJjS0FpY1lOc09PZmtNd2d6aTVYTSUyRlZ5WWF6QXc9IiwiU2VydmVycyI6WyIxMC4yOS4xLjI6MTIzNDUiLCIxMC4yOS4xLjQ6MTIzNDUiLCIxMC4yOS4xLjM6MTIzNDUiLCIxMC4yOS4xLjY6MTIzNDUiLCIxMC4yOS4xLjE6MTIzNDUiLCIxMC4yOS4xLjU6MTIzNDUiXSwiVXVpZHMiOlsiYWZkYzFiYzItYWY2Zi0xMWVjLThkMjEtYTIzNjk1OTM3MTQ0IiwiYWZlMDMxNTYtYWY2Zi0xMWVjLWJjZDUtYTIzNjk1OTM3MTQ0IiwiYWZlNGJlMjEtYWY2Zi0xMWVjLWIxNTYtYTIzNjk1OTM3MTQ0IiwiYWZlOTk0MDUtYWY2Zi0xMWVjLTgyMDItYTIzNjk1OTM3MTQ0IiwiYWZlZTg0YzgtYWY2Zi0xMWVjLWFmNzMtYTIzNjk1OTM3MTQ0IiwiYWZmMmQ1MWYtYWY2Zi0xMWVjLTk1YzEtYTIzNjk1OTM3MTQ0Il19

# 2.Second 有剩下的68000字节的数据
curl -v -XPUT  --data-binary @C:/Users/68725/Desktop/second -H "range: bytes=32000" 10.29.2.1:12345
curl -v -XPUT  --data-binary @C:/Users/68725/Desktop/second -H "range: bytes=32000" 10.29.2.1:12345/temp/eyJOYW1lIjoidGVzdDYiLCJTaXplIjoxMDAwMDAsIkhhc2giOiJrWkxDVzNOUHk2MitNdHJjS0FpY1lOc09PZmtNd2d6aTVYTSUyRlZ5WWF6QXc9IiwiU2VydmVycyI6WyIxMC4yOS4xLjI6MTIzNDUiLCIxMC4yOS4xLjQ6MTIzNDUiLCIxMC4yOS4xLjM6MTIzNDUiLCIxMC4yOS4xLjY6MTIzNDUiLCIxMC4yOS4xLjE6MTIzNDUiLCIxMC4yOS4xLjU6MTIzNDUiXSwiVXVpZHMiOlsiYWZkYzFiYzItYWY2Zi0xMWVjLThkMjEtYTIzNjk1OTM3MTQ0IiwiYWZlMDMxNTYtYWY2Zi0xMWVjLWJjZDUtYTIzNjk1OTM3MTQ0IiwiYWZlNGJlMjEtYWY2Zi0xMWVjLWIxNTYtYTIzNjk1OTM3MTQ0IiwiYWZlOTk0MDUtYWY2Zi0xMWVjLTgyMDItYTIzNjk1OTM3MTQ0IiwiYWZlZTg0YzgtYWY2Zi0xMWVjLWFmNzMtYTIzNjk1OTM3MTQ0IiwiYWZmMmQ1MWYtYWY2Zi0xMWVjLTk1YzEtYTIzNjk1OTM3MTQ0Il19

# 下载这个文件
curl -v 10.29.2.1:12345/objects/test6 -o file
# token
/temp/eyJOYW1lIjoidGVzdDYiLCJTaXplIjoxMDAwMDAsIkhhc2giOiJrWkxDVzNOUHk2MitNdHJjS0FpY1lOc09PZmtNd2d6aTVYTSUyRlZ5WWF6QXc9IiwiU2VydmVycyI6WyIxMC4yOS4xLjI6MTIzNDUiLCIxMC4yOS4xLjQ6MTIzNDUiLCIxMC4yOS4xLjM6MTIzNDUiLCIxMC4yOS4xLjY6MTIzNDUiLCIxMC4yOS4xLjE6MTIzNDUiLCIxMC4yOS4xLjU6MTIzNDUiXSwiVXVpZHMiOlsiYWZkYzFiYzItYWY2Zi0xMWVjLThkMjEtYTIzNjk1OTM3MTQ0IiwiYWZlMDMxNTYtYWY2Zi0xMWVjLWJjZDUtYTIzNjk1OTM3MTQ0IiwiYWZlNGJlMjEtYWY2Zi0xMWVjLWIxNTYtYTIzNjk1OTM3MTQ0IiwiYWZlOTk0MDUtYWY2Zi0xMWVjLTgyMDItYTIzNjk1OTM3MTQ0IiwiYWZlZTg0YzgtYWY2Zi0xMWVjLWFmNzMtYTIzNjk1OTM3MTQ0IiwiYWZmMmQ1MWYtYWY2Zi0xMWVjLTk1YzEtYTIzNjk1OTM3MTQ0Il19


#---------develop7---------#
# 生成100M文件
dd if=/dev/zero of=/tem/file bs=1M count=100
# 获取该文件的sha256 base64编码
openssl dgst -sha256 -binary C:/Users/68725/Desktop/file | base64
IEkqTQ2E+L6xdn9mFiKfhdRMKCe2S9v7Jg7hL6EQng4=
# 上传文件
curl -v -XPUT 10.29.2.1:12345/objects/test7 -H "Digest: SHA-256=IEkqTQ2E+L6xdn9mFiKfhdRMKCe2S9v7Jg7hL6EQng4=" --data-binary @C:/Users/68725/Desktop/file

# 下载文件 可以看到是100M
curl -v  10.29.2.1:12345/objects/test7 -o test7
# 下载压缩文件 发现只下载了99k
curl -v  10.29.2.1:12345/objects/test7 -o test7 -H "accept-encoding: gzip"
# 查看差异
diff file test7


#----------OSS的维护----------#
echo -n "this is object test8 version 1" | openssl dgst -sha256 -binary | base64
2IJQkIth94IVsnPQMrsNxz1oqfrsPo0E2ZmZfJLDZnE=
# 上传第一个版本
curl -v 10.29.2.1:12345/objects/test8 -XPUT -d "this is object test8 version 1" -H "Digest: SHA-256=2IJQkIth94IVsnPQMrsNxz1oqfrsPo0E2ZmZfJLDZnE="
# 上传第二个版本
echo -n "this is object test8 version 2-6" | openssl dgst -sha256 -binary | base64
66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA=

curl -v 10.29.2.1:12345/objects/test8 -XPUT -d "this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl -v 10.29.2.1:12345/objects/test8 -XPUT -d "this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl -v 10.29.2.1:12345/objects/test8 -XPUT -d "this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl -v 10.29.2.1:12345/objects/test8 -XPUT -d "this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="
curl -v 10.29.2.1:12345/objects/test8 -XPUT -d "this is object test8 version 2-6" -H "Digest: SHA-256=66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA="

# 运行deleteOldMetaData工具
export ES_SERVER=localhost:9200
# 发现version1的版本被删除了
go run deleteOldMetaData.go

# 运行deleteOrphanObject工具 此时版本1已经全部被移到Orphan
LISTEN_ADDRESS=10.29.1.1:12345 STORAGE_ROOT=C:/tmp/1 go run deleteOrphanObject.go
LISTEN_ADDRESS=10.29.1.2:12345 STORAGE_ROOT=C:/tmp/2 go run deleteOrphanObject.go
LISTEN_ADDRESS=10.29.1.3:12345 STORAGE_ROOT=C:/tmp/3 go run deleteOrphanObject.go
LISTEN_ADDRESS=10.29.1.4:12345 STORAGE_ROOT=C:/tmp/4 go run deleteOrphanObject.go
LISTEN_ADDRESS=10.29.1.5:12345 STORAGE_ROOT=C:/tmp/5 go run deleteOrphanObject.go
LISTEN_ADDRESS=10.29.1.6:12345 STORAGE_ROOT=C:/tmp/6 go run deleteOrphanObject.go

# 运行objectScanner工具
# 删除一个数据分片
rm C:/tmp/1/objects/66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA=.5.ih70CdjuiOerAQJiRj5Nnha6at+Rz9A6GKemrqmIDD4=
# 修改一个数据分片
echo wxf68725032 > C:/tmp/2/objects/66WuRH0s0albWDZ9nTmjFo9JIqTTXmB6EiRkhTh1zeA=.3.k7Z7BMDLAqtsm+AnQLO0dwSdXat1CnaUgRyE0f9ZgZ0=
# 在节点3上运行数据修复程序 发现节点1上的数据被恢复了   节点2上的数据被修正了
STORAGE_ROOT=C:/tmp/3 go run objectScanner.go




