
pb:
	protoc -I../engula/src/api/  \
	    --go_out=.  \
		--go_opt=Mengula/v1/metadata.proto=pkg/proto \
		--go_opt=Mengula/v1/engula.proto=pkg/proto \
		../engula/src/api/engula/v1/*.proto
	protoc -I../engula/src/api/  \
		--go-grpc_out=. \
		--go-grpc_opt=paths=import \
		--go-grpc_opt=Mengula/v1/metadata.proto=pkg/proto \
		--go-grpc_opt=Mengula/v1/engula.proto=pkg/proto \
		../engula/src/api/engula/v1/engula.proto
