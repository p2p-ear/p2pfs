#include "c_interface.h"

int main() {
    // Initialize the arguments. Note, that we use Go types, defined in c_interface.h.

    //IP - ip of any node on the network
    GoString ip = { "127.0.0.1:9000", 14 };

    //ID - the size of the ring
    GoUint64 ringsz = 10000;

    //fname - the name of the file, same format (last number is the length of the string)
    GoString fname = { "testfile", 8 };
    
    //fcontent - the file's content - as an array
    GoInt8 fcontent[4] = {0, 0, 0, 0};

    //fcontent_slice - the content in required format - GoSlice (last two numbers are size and maxsize of the array respectively)
    GoSlice fcontent_slice = {fcontent, 4, 4};

    
    UploadFile(ip, fname, ringsz, fcontent_slice);
}