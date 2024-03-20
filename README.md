## 项目介绍

作为本人本科毕设***GPU监测系统***后端项目，使用gin搭建web，gorm操作数据库，jwt令牌用户验权，sha1密码加密，向gpt发api请求，自配置seelog日志文件，ssh连接服务器，文件框架仿go-admin，由于同时进行前后端开发，精力有限，仍在陆续修进。

<br><br><br><br>

## 如何运行(需要配合前端一起启动以及golang环境)
**部署项目，在你的代码文件夹里，win+r输入cmd回车**
```
git clone https://github.com/dianayyds/GPU-backend.git
```

**进入gpu-backend文件**

**找到Usersconfig.json并将其修改为你的mysql数据库配置，如果你不知道怎么修改，那么我教你配置,该文件就不用修改了，直接win+r输入cmd回车**
```
mysql -u root -p
(此处输入你的mysql密码)
DROP DATABASE IF EXISTS userinfo;
CREATE DATABASE userinfo;
```

**回到gpu-backend文件夹目录，win+r输入cmd回车，下载依赖**
```
go mod tidy
```

**启动**
```
go run main.go
```
