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
- [Kamva/mgm](https://github.com/Kamva/mgm)
- [9.5 存储密码](https://www.kancloud.cn/kancloud/web-application-with-golang/44198)

### web相关
- [gin获取路径中的参数](https://blog.csdn.net/ma2595162349/article/details/109398069)

### 定时任务
- [cron](https://pkg.go.dev/github.com/robfig/cron#section-readme)

### 环境变量
- [Golang 获取系统环境变量](https://studygolang.com/articles/3387)

### 工具
- [https://robomongo.org/download](https://robomongo.org/download)

### 游戏相关
- [电子游戏术语列表](https://zh.wikipedia.org/wiki/%E9%9B%BB%E5%AD%90%E9%81%8A%E6%88%B2%E8%A1%93%E8%AA%9E%E5%88%97%E8%A1%A8)