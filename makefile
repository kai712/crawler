# 开始爬取
start:
	source ./env.sh; go run ./src/crawler.github.com/main.go
# 安装依赖
install:
	source ./env.sh; cd ./src; glide install; cd ../