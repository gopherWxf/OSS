curl 127.0.0.1:9200/metadata -XDELETE

curl localhost:9200/metadata -XPUT  -H 'Content-Type: application/json' -d'{"mappings":{"objects":{"properties":{"name":{"type":"keyword"},"version":{"type":"integer"},"size":{"type":"integer"},"hash":{"type":"keyword"}}}}}'
