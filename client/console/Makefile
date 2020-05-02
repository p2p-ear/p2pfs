CC=g++

DIR_SRC = ./src/
DIR_INC = ./include/
DIR_BLD = ./build/
DIR_LIB = ../libs/build/

CFLAGS = -c -std=c++1z -O3
LFLAGS = -lduload -L$(DIR_LIB)

all: MDupload MDdownload

preprocess:
	mkdir build

MDupload: MDupload.o
	$(CC) $(DIR_BLD)MDupload.o $(LFLAGS) -o $(DIR_BLD)MDupload

MDupload.o: $(DIR_SRC)console_upload.cpp
	$(CC) $(CFLAGS) $(DIR_SRC)console_upload.cpp -o $(DIR_BLD)MDupload.o

MDdownload: MDdownload.o
	$(CC) $(DIR_BLD)MDdownload.o $(LFLAGS) -o $(DIR_BLD)MDdownload

MDdownload.o: $(DIR_SRC)console_download.cpp
	$(CC) $(CFLAGS) $(DIR_SRC)console_download.cpp -o $(DIR_BLD)MDdownload.o