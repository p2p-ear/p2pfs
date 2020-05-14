#include "libc_interface.h"
#include <assert.h>
#include <stdlib.h>
#include <stdio.h>

GoSlice gen_rand_str(int len);
int testUD(GoString ip, GoUint64 ringsz, GoString fname);
int testUD_RSC(GoString ip, GoUint64 ringsz, GoString fname);

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

    int err = testUD(ip, ringsz, fname);
    if (err != 0) {
        return err;
    }

    err = testUD_RSC(ip, ringsz, fname);
    if (err != 0) {
        return err;
    }

    return 0;
}

GoSlice gen_rand_str(int len) {
    //fcontent - the file's content - as an array
    GoInt8* gobytearr = (GoInt8*) calloc(len, sizeof(GoInt8));
    for (int i = 0; i < len; i++) {
        gobytearr[i] = rand();
    }

    //fcontent_slice - the content in required format - GoSlice (last two numbers are size and maxsize of the array respectively)
    GoSlice goslice = {
        data: gobytearr, 
        len: len, 
        cap: len
    };

    return goslice;
}


int testUD(GoString ip, GoUint64 ringsz, GoString fname) {
    
    const int ARRSZ = 4096;
    GoSlice fcontent_slice = gen_rand_str(ARRSZ);
    
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
        if(fcontent_read[i] != ((GoInt8*)fcontent_slice.data)[i]) {
            fprintf(stderr, "Read and written bytes don't match!");
            return -1;
        }
    }

    DeleteFile(ip, fname, ringsz);

    free(fcontent_read);
    free(fcontent_slice.data);

    remove(fname.p);
    
    return 0;
}

int testUD_RSC(GoString ip, GoUint64 ringsz, GoString fname) {
    
    const int ARRSZ = 4096;
    GoSlice fcontent_slice = gen_rand_str(ARRSZ);
    
    UploadFileRSC(ip, fname, ringsz, fcontent_slice);

    GoSlice buff = {
        //Allocate at least 9 more bytes than nescessary
        data: (GoInt8*) calloc(ARRSZ + 9, sizeof(GoInt8)),
        len: ARRSZ + 9,
        cap: ARRSZ + 9
    };

    GoInt remainingSpace = DownloadFileRSC(ip, fname, ringsz, buff);

    GoInt8* fcontent_read = buff.data;

    if (remainingSpace != 0) {
        fprintf(stderr, "Read and written sizes don't match, empty = %lld", remainingSpace);
        return -1;
    }
    for (int i = 0; i < ARRSZ; i ++) {
        if(fcontent_read[i] != ((GoInt8*)fcontent_slice.data)[i]) {
            fprintf(stderr, "Read and written bytes don't match!");
            return -1;
        }
    }

    DeleteFileRSC(ip, fname, ringsz);

    free(fcontent_read);
    free(fcontent_slice.data);

    remove(fname.p);
    
    return 0;
}