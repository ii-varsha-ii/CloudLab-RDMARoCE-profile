

build:
	#rm -rf data/*.pb.go
	protoc data/rpcmessage.proto --go_out=./ --go-grpc_out=./
	docker build --platform linux/amd64 -t adarshzededa/rdma-redis-agent .

push:
	docker push adarshzededa/rdma-redis-agent:latest