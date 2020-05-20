SERVERPATH = storagePeer
CLIENTPATH = client

shared:
	cd ${SERVERPATH}/ ; $(MAKE)
	cp ${SERVERPATH}/c_interface/c_interface.h ${CLIENTPATH}/libs/include
