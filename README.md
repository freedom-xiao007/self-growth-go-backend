# 自我生成--后台服务
***

## HTTPS证书相关
### Windows10
```shell
& 'C:\Program Files\OpenSSL-Win64\bin\openssl.exe' genrsa -des3 -out server.key 2048
& 'C:\Program Files\OpenSSL-Win64\bin\openssl.exe' req -new -key server.key -out server.csr
& 'C:\Program Files\OpenSSL-Win64\bin\openssl.exe' rsa -in server.key -out server_no_passwd.key
& 'C:\Program Files\OpenSSL-Win64\bin\openssl.exe' x509 -req -days 365 -in server.csr -signkey server_no_passwd.key -out server.crt
```

## 参考链接
- [Go语言项目开发实战](https://time.geekbang.org/column/article/381392)