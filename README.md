# Environment for launch system metrics
## Prerequisites
export PATH=$PATH:/usr/local/go/bin__
export PATH="$PATH:$(go env GOPATH)/bin"__
go mod init project__
go mod tidy__
apt install -y protobuf-compiler__
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest__
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest__
protoc --go_out=. --go-grpc_out=. proto/agent.proto__
# In folder object calls where place main.go do:
go run . 