
gen-pb:
	protoc --go_out=. --go-grpc_out=. broker.proto