
## ⚙ 配置详情

1. 配置文件说明
```yml
settings:
  application:  
    # 项目启动环境            
    mode: dev  # dev开发环境 test测试环境 prod线上环境；
    host: 0.0.0.0  # 主机ip 或者域名，默认0.0.0.0
    # 服务名称
    name: admin   
    # 服务端口
    port: 8000   
    readtimeout: 1   
    writertimeout: 2 
  log:
    # 日志文件存放路径
    dir: temp/logs
  jwt:
    # JWT加密字符串
    secret: admin
    # 过期时间单位：秒
    timeout: 3600
  database:
    # 数据库名称
    name: dbname 
    # 数据库类型
    dbtype: mysql    
    # 数据库地址
    host: 127.0.0.1  
    # 数据库密码
    password: password  
    # 数据库端口
    port: 3306       
    # 数据库用户名
    username: root   
```

2. 文件路径  config/settings.yml


## 📦 本地开发

### 首次启动说明

```bash

# 编译项目
go build

# 修改配置
vi ./config/setting.yml 

# 1. 配置文件中修改数据库信息 
# 注意: settings.database 下对应的配置数据
# 2. 确认log路径

```

### 初始化数据库，以及服务启动
```
# 首次配置需要初始化数据库资源信息
./main.exe init -c config/settings.yml -m dev


# 启动项目，也可以用IDE进行调试
./main.exe server -c config/settings.yml -p 8000 -m dev

windows终止进程 

> netstat -aon | findstr "8000"
> tasklist|findstr "9316"
> taskkill /IM   main.exe  /F

```

### 文档生成
```bash
swag init

# 如果没有swag命令 go get安装一下即可
go get -u github.com/swaggo/swag/cmd/swag
```

### 交叉编译
```bash
env GOOS=windows GOARCH=amd64 go build main.go

# or

env GOOS=linux GOARCH=amd64 go build main.go
```