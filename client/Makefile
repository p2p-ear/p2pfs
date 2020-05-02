all: making moving installing

remove:
	cd scripts ; bash ./uninstall.sh

making:
	cd libs ; $(MAKE) ; cd ../console ; $(MAKE) 

moving:
	cp ./console/build/MDupload ./build ; cp ./console/build/MDdownload ./build

installing:
	cd scripts ; bash ./install.sh


