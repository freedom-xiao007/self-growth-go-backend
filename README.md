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

## 构建与部署相关
### Docker镜像工具
```shell
go get -u github.com/tal-tech/go-zero/tools/goctl
cd .\internal\growth_record
goctl docker -go run.go
docker build -t controller_api:v1 -f internal/growth_record/Dockerfile .

cd .\internal\game_text_job\
goctl docker -go main.go
docker build -t game_text_job:v1 -f internal/game_text_job/Dockerfile .
```

### Portainer
```yaml
version: '2'
services:
  # Go 后台Api服务
  controller4g:
    image: controller_api:v1
    ports:
      - 8081:8080
    environment:
      - mongo_user=user
      - mongo_password=password
      - mongo_host=127.0.0.1
      - mongo_port=27017
      
  # Go 文本游戏服务后台
  gameTextJob:
    image: game_text_job:v1
    environment:
      - mongo_user=user
      - mongo_password=password
      - mongo_host=127.0.0.1
      - mongo_port=27017
      
networks:
  default:
    external:
      name: self_growth 
```

## 参考链接
- [Go语言项目开发实战](https://time.geekbang.org/column/article/381392)
- [Kamva/mgm](https://github.com/Kamva/mgm)
- [9.5 存储密码](https://www.kancloud.cn/kancloud/web-application-with-golang/44198)
- [汉语拼音转换工具 Go 版。](https://pkg.go.dev/github.com/mozillazg/go-pinyin#section-readme)
- [Go生成随机数](https://blog.csdn.net/u011304970/article/details/72721747)

### web相关
- [gin获取路径中的参数](https://blog.csdn.net/ma2595162349/article/details/109398069)

### 定时任务
- [cron](https://pkg.go.dev/github.com/robfig/cron#section-readme)

### 环境变量
- [Golang 获取系统环境变量](https://studygolang.com/articles/3387)

### 工具
- [https://robomongo.org/download](https://robomongo.org/download)
- [最简单的Go Dockerfile编写姿势，没有之一！](https://segmentfault.com/a/1190000038437935)

### 游戏相关
- [电子游戏术语列表](https://zh.wikipedia.org/wiki/%E9%9B%BB%E5%AD%90%E9%81%8A%E6%88%B2%E8%A1%93%E8%AA%9E%E5%88%97%E8%A1%A8)
- [上古神祇 与 神话人物](https://zh.wikipedia.org/wiki/%E4%B8%AD%E5%9B%BD%E7%A5%9E%E8%AF%9D%E4%BA%BA%E7%89%A9%E5%88%97%E8%A1%A8)
- [三毒、八苦，您知道多少？】](http://blog.sina.com.cn/s/blog_62d64c350100k1pt.html)
- [电子游戏术语列表](https://zh.wikipedia.org/wiki/%E9%9B%BB%E5%AD%90%E9%81%8A%E6%88%B2%E8%A1%93%E8%AA%9E%E5%88%97%E8%A1%A8)