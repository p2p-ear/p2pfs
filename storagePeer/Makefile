BIN_PATH    = bin
C_INTERFACE = c_interface

#===========#
# Commands
#===========#

.ONESHELL:

build:
	go build -o ${BIN_PATH}/peer .

gen_proto_ring:
	cd src
	protoc -I dht/ \
		-Idht \
		--go_out=plugins=grpc:dht \
		dht/ring.proto

gen_proto_peer:
	cd src
	protoc -I peer  \
		--go_out=plugins=grpc:peer \
		peer/peer.proto

gen_c_interface:
	cd ${C_INTERFACE}
	go build -o $(addsuffix .so, ${C_INTERFACE}) -buildmode=c-shared $(addsuffix .go, ${C_INTERFACE})
	mv $(addsuffix .so, ${C_INTERFACE}) ../bin

gen_c_test:
	cd ${C_INTERFACE}
	gcc test.c -o ../bin/c_test ../bin/${C_INTERFACE}.so

test:
	cd src
	go test ./*/

run:
	${BIN_PATH}/peer
