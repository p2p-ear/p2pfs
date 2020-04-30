#include "c_interface.h"
#include <assert.h>
#include <stdlib.h>
#include <stdio.h>

int main() {
    // Initialize the arguments. Note, that we use Go types, defined in c_interface.h.

    //IP - ip of any node on the network
    GoString ip = { "127.0.0.1:9000", 14 };

    //ID - the size of the ring
    GoUint64 ringsz = 10000;

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