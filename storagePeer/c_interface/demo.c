#include "c_interface.h"

int main() {
    GoString ip = { "127.0.0.1:9000", 14 };
    GoUint64 id = 0;
    GoString fname = { "testfile", 8 };
    
    GoInt8 fcontent[4] = {0, 0, 0, 0};
    GoSlice fcontent_slice = {fcontent, 4, 4};
    UploadFile(ip, id, fname, fcontent_slice);
}