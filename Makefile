.PHONY: proto

proto:
	protoc \
		--proto_path=proto \
		--go_out=server/proto --go_opt=paths=source_relative \
		--csharp_out=client/proto \
		proto/game.proto
