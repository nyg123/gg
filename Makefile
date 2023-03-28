# 启动本地环境
.PHONY: dev
dev:
	go build main.go && ./main -c=./config/dev.yaml -alsologtostderr=true
