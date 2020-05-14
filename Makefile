SERVERPATH = storagePeer
CLIENTPATH = client

shared:
	cd ${SERVERPATH}/ ; $(MAKE) gen_c_interface
	mkdir -p ${CLIENTPATH}/libs/build
	cp ${SERVERPATH}/bin/libc_interface.so ${CLIENTPATH}/libs/build
	cp ${SERVERPATH}/c_interface/c_interface.h ${CLIENTPATH}/libs/include
