# NOTE: this is untested
python3 -m grpc_tools.protoc -I ./results/proto --grpc_python_out=./results/server/generated --mypy_out=./results/server/generated ./results/proto/results.proto

protoc --go_out=client/proto --go_opt=paths=source_relative --go-grpc_out=client/proto --go-grpc_opt=paths=source_relative demo.proto