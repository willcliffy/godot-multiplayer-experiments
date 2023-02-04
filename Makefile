.PHONY: proto

proto: # --csharp_out=client/proto
	protoc --proto_path=proto --go_out=server/proto --go_opt=paths=source_relative proto/game.proto
