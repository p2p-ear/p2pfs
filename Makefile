SERVERPATH = storagePeer
CLIENTPATH = client

shared:
	cd ${SERVERPATH}/ ; $(MAKE) gen_c_interface
<<<<<<< HEAD
	cp ${SERVERPATH}/bin/c_interface.so ${CLIENTPATH}/libs/build
	cp ${SERVERPATH}/c_interface/c_interface.h ${CLIENTPATH}/libs/include
=======
	mkdir -p ${CLIENTPATH}/libs/build
	cp ${SERVERPATH}/bin/libc_interface.so ${CLIENTPATH}/libs/build
	cp ${SERVERPATH}/c_interface/c_interface.h ${CLIENTPATH}/libs/include
>>>>>>> f4d5f55a5dc2a0e31e99b848d8a2c26d26c97849
