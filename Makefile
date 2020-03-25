#===========#
# Constants
#===========#

BIN_PATH  = bin
CODE_PATH = storagePeer

#===========#
# Commands
#===========#

.ONESHELL:

build:
	cd ${CODE_PATH}
	go build -o ${BIN_PATH}/peer .

gen_proto_ring:
	cd ${CODE_PATH}/src
	protoc -I dht/ \
    -Idht \
    --go_out=plugins=grpc:dht \
    dht/ring.proto

run:
	cd ${CODE_PATH}
	${BIN_PATH}/peer
