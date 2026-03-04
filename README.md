# Environment for launch system metrics
## Prerequisites
export PATH=$PATH:/usr/local/go/bin <br />
export PATH="$PATH:$(go env GOPATH)/bin" <br />
go mod init project <br />
go mod tidy <br />
apt install -y protobuf-compiler <br />
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest <br />
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest <br />
protoc --go_out=. --go-grpc_out=. proto/agent.proto <br />
# In folder object calls where place main.go do:
go run . 