#===========#
# Constants
#===========#

BIN_PATH    = bin
CODE_PATH   = storagePeer
C_INTERFACE = c_interface

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

gen_proto_peer:
	cd ${CODE_PATH}/src
	protoc -I peer  \
		--go_out=plugins=grpc:peer \
		peer/peer.proto

gen_c_interface:
	cd ${CODE_PATH}/${C_INTERFACE}
	go build -o $(addsuffix .so, ${C_INTERFACE}) -buildmode=c-shared $(addsuffix .go, ${C_INTERFACE})
	cp $(addsuffix .so, ${C_INTERFACE}) ../bin

gen_c_test:
	cd ${CODE_PATH}/${C_INTERFACE}
	gcc test.c -o ../bin/c_test ../bin/${C_INTERFACE}.so

test:
	cd ${CODE_PATH}/src
	go test ./*/

run:
	cd ${CODE_PATH}
	${BIN_PATH}/peer
