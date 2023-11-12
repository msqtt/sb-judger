.DEFAULT_GOAL := sb-judger
untar:
	tar zxvf rootfs.tgz
sandbox: cmd/sandbox/main.go
	go build -ldflags="-w -s" -o sandbox cmd/sandbox/main.go
sb-judger: sandbox cmd/judger/main.go
	go build -ldflags="-w -s" -o sb-judger cmd/judger/main.go
gen:
	protoc --proto_path=api/protos/v1 \
	--go_out=api/pb/v1 --go_opt=paths=source_relative \
	--go-grpc_out=api/pb/v1 --go-grpc_opt=paths=source_relative  \
	--grpc-gateway_out=api/pb/v1 --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=api/openapi/v1 \
		api/protos/v1/sandbox/*.proto \
		api/protos/v1/judger/*.proto
dev:
	go run cmd/judger/main.go
