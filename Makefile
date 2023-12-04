# deploy cmd
sb-judger: sandbox cmd/judger/main.go
	go build -ldflags="-w -s" -o sb-judger cmd/judger/main.go

sandbox: cmd/sandbox/main.go
	go build -ldflags="-w -s" -o sandbox cmd/sandbox/main.go

clean:
	rm sandbox sb-judger

docker: build/Dockerfile rootfs
	 docker build -f build/Dockerfile -t msqt/sb-judger:0.1.0 .

rootfs: build/tarball.Dockerfile
	docker build -f build/tarball.Dockerfile -t msqt/rootfs-tarball:0.1.0 .
	mkdir rootfs
	docker create msqt/rootfs-tarball:dev | xargs docker export | tar -C rootfs -xf -
	rm -rf ./rootfs/.dockerenv ./rootfs/var/* ./rootfs/dev/*

# develop cmd
gen:
	protoc --proto_path=api/protos/v1 \
	--go_out=api/pb/v1 --go_opt=paths=source_relative \
	--go-grpc_out=api/pb/v1 --go-grpc_opt=paths=source_relative  \
	--grpc-gateway_out=api/pb/v1 --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=api/openapi/v1 \
		api/protos/v1/sandbox/*.proto \
		api/protos/v1/judger/*.proto

dev: sandbox
	go run cmd/judger/main.go

.PHONY: docker gen dev
