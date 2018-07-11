自定义服务器自签名证书

```

# 生成私钥
openssl genrsa -out server.key 2048

# 生成公钥。 里面的 serverName 填成自己的。 在go里是serverName
openssl req -x509 -new -nodes -key server.key -subj "/CN=severName" -days 5000 -out server.pem

```# grpcDemo
