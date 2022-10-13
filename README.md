# 电商系统-微服务 mall_srvs

##### 一、用户微服务（user_srv）：  

+ 1 用户服务 
  + 添加用户
  + 更新用户
  + 查询用户——通过id
  + 查询用户——通过mobile
  + 用户列表查看  
  + 检查密码
+ 用户权限验证
  + jwt

环境:

+ sdk:go version go1.18.4 windows/amd64

+ 开发工具：goland2021.3.3
+ 技术栈 go、gin、gorm、grpc、consul、nacos、docker
+ 数据库：mysql8、redis5.0


####  二、安装下载

2.1、下载:

grpc服务客户端

>  git clone仓库代码：git clone git@github.com:xieHuiHuang/mall_srvs.git
> go mod tidy

2.2 相关服务配置

grpc服务

> cd proto
>
> protoc -I . user.proto --go_out=plugins=grpc:.
> 
+ 2.3 数据库服务
  + 下载镜像：
  > ```shell
  docker pull mysql:8.0.22
  > ```
  + 查看镜像：
  > ```shell
  docker images
  > ```
  
  + 通过镜像启动：
  > ```shell
  > docker run -p 3306:3306 --name mysql -v $PWD/data/developer/mysql/conf:/etc/mysql/conf.d -v  $PWD/data/developer/mysql/logs:/logs -v  $PWD/data/developer/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:8.0.22
  > ```
  + 重启容器：docker restart mysql
  + 查看镜像：docker ps -a
  
+ 2.4 consul注册
  + 拉取镜像：docker pull consul:latest
  + 安装
  > ```shell
  > docker run -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0
  > docker container update --restart=always 容器id 
  > ```
  + 浏览器访问：http://127.0.0.1:8500

+ 2.4 配置中心nacos
+ 下载
> ```shell
> docker run --name nacos-standalone -e MODE=standalone -e JVM_XMS=512m -e JVM_XMX=512m -e JVM_XMN=256m -p 8848:8848 -d nacos/nacos-server:latest
> ```
+ 访问： http://127.0.0.1:8848/nacos
+ 登录：账号密码都是nacos
+ 根据config.yaml对nacos进行相关配置，只需要配置命名空间及配置列表，配置信息如下
```json
{
  "name": "mall_srvs",
  "mysql": {
    "host": "127.0.0.1",
    "port": 3306,
    "user": "root",
    "password": "123456",
    "db": "mall_user_srv"
  },
  "log": {
    "level": "debug",
    "filename": "web_app_log.log",
    "max_size": 200,
    "max_age": 30,
    "max_backups": 7
  },
  "consul": {
    "host": "127.0.0.1",
    "port": 8500,
      # 更换本机的ipv4地址
    "serverhost": "192.168.31.101"
  }
}
```

#### 三 运行：

> go  run main.go



