#!/usr/bin/env sh

Country="CN"
State="GZ"
Location="GZ"
Organization="ronething"
Organizational="ronething"
CommonName="localhost"

echo "生成 ca 密钥"
openssl genrsa -out ca.key 2048

echo "生成 ca 公钥"
# openssl req -new -x509 -days 3650 -key ca.key -out ca.pem
# 无交互生成公钥
openssl req -new -x509 -key ca.key -out ca.pem -days 3650 -subj "/C=$Country/ST=$State/L=$Location/O=$Organization/OU=$Organizational/CN=$CommonName"

echo "生成服务端 SAN 证书"
openssl genpkey -algorithm RSA -out server.key
openssl req -new -nodes -key server.key -out server.csr -days 3650 -subj "/C=$Country/O=$Organization/OU=$Organizational/CN=$CommonName" -config ./openssl.cnf -extensions v3_req
openssl x509 -req -days 3650 -in server.csr -out server.pem -CA ca.pem -CAkey ca.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req
