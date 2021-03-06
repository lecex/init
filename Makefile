.PHONY: git
git:
	git add .
	git commit -m"自动提交 git 代码"
	git push
.PHONY: tag
tag:
	git push --tags
.PHONY: rpc
rpc:
	micro api  --handler=rpc  --namespace=go.micro.api --address=:8080
.PHONY: api
api:
	micro api  --handler=api  --namespace=go.micro.api --address=:8081

.PHONY: proto
proto:
	protoc -I . --micro_out=. --gogofaster_out=. proto/health/health.proto

.PHONY: docker
docker:
	docker build -f Dockerfile  -t init .
.PHONY: run
run:
	go run main.go
test:
	go test main_test.go -test.v
t:
	git tag -d v1.5.2
	git push origin :refs/tags/v1.5.2
	git tag v1.5.2
	git push --tags