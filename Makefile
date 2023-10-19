gen:
	protoc --proto_path=api/proto \
				 --go_out=api/proto/pb --go_opt=paths=source_relative \
		api/proto/*.proto
