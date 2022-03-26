
# 计算内容的哈希值
echo -n "this is object test3 version 1" | openssl dgst -sha256 -binary | base64
#p3NoRweZQOBNCrQa0QuCVKV2RiDWEst1/G+O8dxafq0=
# put
curl -v 10.29.2.1:12345/objects/test3 -XPUT -d "this is object test3 version 1" -H "Digest: SHA-256=p3NoRweZQOBNCrQa0QuCVKV2RiDWEst1/G+O8dxafq0="
# 计算内容的哈希值
echo -n "this is object test3 version 2" | openssl dgst -sha256 -binary | base64
#cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo=
# put
curl -v 10.29.2.1:12345/objects/test3 -XPUT -d "this is object test3 version 2" -H "Digest: SHA-256=cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo="
# 删除
curl -v 10.29.2.1:12345/objects/test3 -XDELETE
# 查看所有版本
curl -v 10.29.2.1:12345/versions/
# 查看指定版本
curl -v 10.29.2.1:12345/objects/test3?version=1