
GOTOOLS = \
	github.com/golang/dep/cmd/dep \
	gopkg.in/alecthomas/gometalinter.v2 \
	github.com/golang/protobuf/proto \
	github.com/golang/protobuf/ptypes/struct \
	google.golang.org/grpc \
	github.com/gogo/protobuf/proto \
	github.com/gogo/protobuf/jsonpb \
	github.com/gogo/protobuf/protoc-gen-gogo \
	github.com/gogo/protobuf/gogoproto \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
	github.com/lib/pq

PROTOPATH = -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf:. -I=${GOPATH}/src/github.com/gallactic/gallactic/rpc/grpc/proto3 -I=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
#--proto_path=${GOPATH}/src:${GOPATH}/src/github.com/gogo/protobuf/protobuf:.
HUBBLE = ${GOPATH}/src/github.com/BurrowBlocks

########################################
### make all
all: tools deps build

########################################
### Tools & dependencies
deps:

	dep ensure

tools:

	go get $(GOTOOLS)
	@gometalinter.v2 --install

########################################
### Protobuf
#proto:

#	--protoc $(PROTOPATH) --gogo_out=plugins=grpc:$(HUBBLE) ./proto3/blockchain.proto

########################################
### Formatting, linting, and vetting
fmt:
	@go fmt ./...

########################################
### building
build:
	@go build main.go

run:
	@go run main.go

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: tools deps
.PHONY: build 
.PHONY: fmt metalinter