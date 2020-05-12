C_INTERFACE = c_interface

#===========#
# Commands
#===========#

.ONESHELL:

build: $(addprefix src/,peer/peer.go dht/ring.go dht/ring_struct.go)
	go build -o bin/peer .

$(addprefix src/peer/,peer.go peer_test.go): src/peer/peer.proto
$(addprefix src/dht/,ring.go ring_struct.go ring_test.go): src/dht/ring.proto

src/dht/ring.proto:
	cd src
	protoc -I dht/ \
		-Idht \
		--go_out=plugins=grpc:dht \
		dht/ring.proto

src/peer/peer.proto:
	cd src
	protoc -I peer  \
		--go_out=plugins=grpc:peer \
		peer/peer.proto

c_interface/c_interface.go: mkbin
	cd c_interface
	go build -o libc_interface.so -buildmode=c-shared c_interface.go
	mv libc_interface.so ../bin

c_interface/test.c: mkbin c_interface/c_interface.go
	cd c_interface
	gcc test.c -o ../bin/c_test ../bin/libc_interface.so

test: c_interface/test.c $(addprefix src/,peer/peer_test.go dht/ring_test.go)
	cd src
	go test ./*/

mkbin:
	mkdir -p bin/

clean:
	rm -rf bin/*

gen_c_interface: c_interface/c_interface.go