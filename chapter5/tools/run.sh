echo -n "this is object test3 version 1" | openssl dgst -sha256 -binary | base64
#p3NoRweZQOBNCrQa0QuCVKV2RiDWEst1/G+O8dxafq0=
curl -v 10.29.2.1:12345/objects/test3 -XPUT -d "this is object test3 version 1" -H "Digest: SHA-256=p3NoRweZQOBNCrQa0QuCVKV2RiDWEst1/G+O8dxafq0="

echo -n "this is object test3 version 2" | openssl dgst -sha256 -binary | base64
#cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo=
curl -v 10.29.2.1:12345/objects/test3 -XPUT -d "this is object test3 version 2" -H "Digest: SHA-256=cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo="

curl -v 10.29.2.1:12345/objects/test3

curl -v 10.29.2.1:12345/objects/test3 -XDELETE

curl -v 10.29.2.1:12345/locate/p3NoRweZQOBNCrQa0QuCVKV2RiDWEst1/G+O8dxafq0=

curl -v 10.29.2.1:12345/objects/test4 -XPUT -d "this is object test4" -H "Digest: SHA-256=qxn98QM9ZPosMkBIKOwQUTI+5s2a4sDNaBBlTT5jLhw="

curl -v 10.29.2.1:12345/objects/test4


curl -v 10.29.2.1:12345/objects/test5 -XPUT -d "this is object test5" -H "Digest: SHA-256=B494C1vj+98Y+PTGRiNqWu7gRgWQwiHnEofa47sN6mk="
echo -n "this is object test4" | openssl dgst -sha256 -binary | base64

curl -v 10.29.2.1:12345/objects/test5 -XDELETE

curl -v 10.29.2.1:12345/locate/B494C1vj+98Y+PTGRiNqWu7gRgWQwiHnEofa47sN6mk=