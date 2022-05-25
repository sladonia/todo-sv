pb:
	protoc -I=proto --go_out=pkg --go-grpc_out=pkg proto/todo/*.proto

.PHONY: pb
