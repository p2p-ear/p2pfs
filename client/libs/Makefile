CFLAGS_LIB =-c -std=c++1z -O3 -fPIC
LFLAGS_LIB = -shared -lstdc++fs

CC=g++

DIR_SRC = ./src/
DIR_INC = ./include/
DIR_BLD = ./build/

all: libduload.so

preprocess:
	mkdir build

clean:
	rm -r build

libduload.so: upload.o download.o common.o 
	$(CC) $(LFLAGS_LIB) $(DIR_BLD)upload.o $(DIR_BLD)download.o $(DIR_BLD)common.o -lc_interface -o $(DIR_BLD)libduload.so 
upload.o: $(DIR_INC)duload.h $(DIR_SRC)upload.cpp 
	$(CC) $(CFLAGS_LIB) $(DIR_SRC)upload.cpp -o $(DIR_BLD)upload.o

download.o: $(DIR_INC)duload.h $(DIR_SRC)download.cpp
	$(CC) $(CFLAGS_LIB) $(DIR_SRC)download.cpp -o $(DIR_BLD)download.o

common.o: $(DIR_INC)duload.h $(DIR_SRC)common.cpp
	$(CC) $(CFLAGS_LIB) $(DIR_SRC)common.cpp -o $(DIR_BLD)common.o
