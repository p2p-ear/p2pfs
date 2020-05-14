#include "libc_interface.h"
#include <assert.h>
#include <stdlib.h>
#include <stdio.h>

int main(int argc, char* argv[]) {
    // Initialize the arguments. Note, that we use Go types, defined in c_interface.h.
    if (argc < 3) {
        fprintf(stderr, "Not enough arguments for c_test: %d < 3", argc);
        return -1;
    }
    
    //IP - ip of any node on the network
    const char* ip_string = argv[1];
    GoString ip = { ip_string, 14 };

    //ID - the size of the ring
    GoUint64 ringsz = atoi(argv[2]);

    //fname - the name of the file, same format (last number is the length of the string)
    GoString fname = { "testfile", 8 };
    
    //fcontent - the file's content - as an array
    const int ARRSZ = 1024;
    GoInt8* fcontent = (GoInt8*) calloc(ARRSZ, sizeof(GoInt8));
    for (int i = 0; i < ARRSZ; i++) {
        fcontent[i] = rand();
    }

    //fcontent_slice - the content in required format - GoSlice (last two numbers are size and maxsize of the array respectively)
    GoSlice fcontent_slice = {
        data: fcontent, 
        len: ARRSZ, 
        cap: ARRSZ
    };

    
    UploadFile(ip, fname, ringsz, fcontent_slice);

    GoSlice buff = {
        data: (GoInt8*) calloc(ARRSZ, sizeof(GoInt8)),
        len: ARRSZ,
        cap: ARRSZ
    };

    GoInt remainingSpace = DownloadFile(ip, fname, ringsz, buff);

    GoInt8* fcontent_read = buff.data;

    if (remainingSpace != 0) {
        fprintf(stderr, "Read and written sizes don't match, empty = %lld", remainingSpace);
        return -1;
    }
    for (int i = 0; i < ARRSZ; i ++) {
        if(fcontent_read[i] != fcontent[i]) {
            fprintf(stderr, "Read and written bytes don't match!");
            return -1;
        }
    }

    free(fcontent_read);
    free(fcontent);

    remove(fname.p);
}