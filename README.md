crawler 是一个基于Go语言编写的高并发的爬虫DEMO，有常见爬取的三种方式，包括“解析页面”、“解析异步ajax数据”、“SPA应用爬取方式”代码实现,并展示到前台。

![](https://github.com/kai712/crawler/tree/master/static/crawler.gif)


### 本项目有以下几个核心功能块

- [x] 基于Golang的三种常见爬取方式实现

- [x] 前端数据展示以及elasticsearch全文检索 （进行中。。）

- [x] 并发日志采集，并实时推送到前端，动态展示爬取详情 （进行中。。）

### 运行项目

将项目clone到本地
```
git clone git@github.com:kai712/crawler.git
```

打开env.sh文件，配置自己的gopath
```
// 项目本地路径
export GOPATH=$HOME/ck/crawler
```

打开main.go文件，配置自己本地mysql
```
// mysql配置
const (
	username = ""
	password = ""
  // 本地数据库名
	dbname   = ""
)
```

运行项目
```
make start
```