echo -n "this is object test3 version 1" | openssl dgst -sha256 -binary | base64
#p3NoRweZQOBNCrQa0QuCVKV2RiDWEst1/G+O8dxafq0=
curl -v 10.29.2.1:12345/objects/test3 -XPUT -d "this is object test3 version 1" -H "Digest: SHA-256=p3NoRweZQOBNCrQa0QuCVKV2RiDWEst1/G+O8dxafq0="

echo -n "this is object test3 version 2" | openssl dgst -sha256 -binary | base64
#cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo=
curl -v 10.29.2.1:12345/objects/test3 -XPUT -d "this is object test3 version 2" -H "Digest: SHA-256=cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo="

curl -v 10.29.2.1:12345/objects/test3

curl -v 10.29.2.1:12345/objects/test3 -XDELETE

curl -v 10.29.2.1:12345/locate/p3NoRweZQOBNCrQa0QuCVKV2RiDWEst1/G+O8dxafq0=

curl -v 10.29.2.1:12345/objects/test4 -XPUT -d "this is object test3 version 1" -H "Digest: SHA-256=p3NoRweZQOBNCrQa0QuCVKV2RiDWEst1/G+O8dxafq0="

curl -v 10.29.2.1:12345/objects/test4