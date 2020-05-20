#include "libc_interface.h"
#include <assert.h>
#include <stdlib.h>
#include <stdio.h>

GoSlice gen_rand_str(int len);
int testUD_RSC(GoString ip, GoUint64 ringsz, GoString fname, GoString rJWT, GoString wJWT, GoString dJWT);

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

    char* secret = "qwertyuiopasdfghjklzxcvbnm123456";
    char* read_JWT = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdGZpbGUiLCJzaXplIjo1MTIsImFjdCI6MH0.S01aGN53BOds6R8JhOfJVWnnJXk8jMM78DqJFaGAMucyKwEhvapx7UzkDqulyU9qrGJrHFgJrZrWzCsydeCtiQ";
    char* writ_JWT = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdGZpbGUiLCJzaXplIjo1MTIsImFjdCI6MX0.w8Npnxi5aLhe9mtZ1pVs9nye_JT6EbBTtm0cLUp7MaWyoHU9wMk4WZoBxRou3KscJaKKFhqM90pzecshsJK_jw";
    char* dele_JWT = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdGZpbGUiLCJzaXplIjo1MTIsImFjdCI6Mn0.az56Tnh8TNNF-JQa42tEqQuAuGx8xIlm9HSUU6QFuTox7jliiaCvzZdxszu7ZuFFkzTFkmlVHSV5MNJGUHvqjg";

    GoString rJWT = {
        p: read_JWT,
        n: 176
    };

    GoString wJWT = {
        p: writ_JWT,
        n: 176        
    };

    GoString dJWT = {
        p: dele_JWT,
        n: 176
    };

    int err = testUD_RSC(ip, ringsz, fname, rJWT, wJWT, dJWT);
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

int testUD_RSC(GoString ip, GoUint64 ringsz, GoString fname, GoString rJWT, GoString wJWT, GoString dJWT) {
    
    const int ARRSZ = 4096;
    GoSlice fcontent_slice = gen_rand_str(ARRSZ);
    
    UploadFileRSC(ip, fname, ringsz, fcontent_slice, wJWT);

    GoSlice buff = {
        //Allocate at least 9 more bytes than nescessary
        data: (GoInt8*) calloc(ARRSZ + 9, sizeof(GoInt8)),
        len: ARRSZ + 9,
        cap: ARRSZ + 9
    };

    GoInt remainingSpace = DownloadFileRSC(ip, fname, ringsz, buff, rJWT);

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

    DeleteFileRSC(ip, fname, ringsz, dJWT);

    free(fcontent_read);
    free(fcontent_slice.data);

    remove(fname.p);
    
    return 0;
}