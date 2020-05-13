SERVERPATH = storagePeer
CLIENTPATH = client

shared:
	cd ${SERVERPATH}/ ; $(MAKE) gen_c_interface
	cp ${SERVERPATH}/bin/c_interface.so ${CLIENTPATH}/libs/build
	cp ${SERVERPATH}/c_interface/c_interface.h ${CLIENTPATH}/libs/include
