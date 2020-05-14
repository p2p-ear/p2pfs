C_INTERFACE = c_interface

PEER = src/peer
RING = src/dht

# All non-generated .go files
PEER_GO = $(filter-out ${PEER}/peer.pb.go,$(wildcard ${PEER}/*.go))
RING_GO = $(filter-out ${RING}/ring.pb.go,$(wildcard ${RING}/*.go))

# Test files
PEER_TEST = ${PEER}/peer_test.go
RING_TEST = ${RING}/ring_test.go
TEST = ${PEER_TEST}${RING_TEST}

# Source files (non-test)
PEER_SRC = ${filter-out ${PEER_TEST},${PEER_GO}}
RING_SRC = ${filter-out ${RING_TEST},${RING_GO}}
#===========#
# Commands
#===========#

.ONESHELL:

build: ${PEER_GO} ${RING_GO}
	go build -o bin/peer .

# Peer
${PEER_SRC}: ${PEER}/peer.pb.go
${PEER}/peer.pb.go: ${PEER}/peer.proto
	cd src
	protoc -I peer  \
		--go_out=plugins=grpc:peer \
		peer/peer.proto
${PEER_TEST}: ${PEER_SRC} c_interface/test.c
	go test ./${PEER}

# Ring
${RING_SRC}: ${RING}/ring.pb.go
${RING}/ring.pb.go: ${RING}/ring.proto
	cd src
	protoc -I dht/ \
		-Idht \
		--go_out=plugins=grpc:dht \
		dht/ring.proto
${RING_TEST}: ${RING_SRC}
	go test ./${RING}

c_interface/c_interface.go: mkbin ${PEER_SRC}
	cd c_interface
	go build -o libc_interface.so -buildmode=c-shared c_interface.go
	mv libc_interface.so ../bin

c_interface/test.c: mkbin c_interface/c_interface.go
	cd c_interface
	gcc test.c -o ../bin/c_test ../bin/libc_interface.so

mkbin:
	mkdir -p bin/

clean:
	rm -rf bin/*

gen_c_interface: c_interface/c_interface.go