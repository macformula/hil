# NOTE: this is untested
python3 -m grpc_tools.protoc -I ./proto --python_out=./server/generated --pyi_out=./server/generated --grpc_python_out=./server/generated ./proto/results.proto

protoc --go_out=./client/generated --go_opt=paths=source_relative --go-grpc_out=./client/generated --go-grpc_opt=paths=source_relative