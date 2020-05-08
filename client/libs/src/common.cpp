#include "../include/duload.h"

//finding out the position of last slash
int LastSlash(const std::string& path) {
    int result = 0;
    for (int i = 0; i < path.length(); i++) {
        if (path[i] == '/') {
            result = i;
        }
        if (path[i] == '?' || path[i] == '*') {
            break;
        }
    }
    return result;
}

int RemoveFile(const std::string& filename) {
    std::vector<char*> rmargs;
    std::string com("rm");
    std::string fl("-rf");
    rmargs.push_back((char*)com.c_str());
    rmargs.push_back((char*)fl.c_str());
    rmargs.push_back((char*)filename.c_str());
    rmargs.push_back(nullptr);
    int pid1 = fork();
    if (pid1 == -1) {
        return 0;
    } 
    if (pid1 == 0) { //child
        //std::cout << filename << "\n";
        execvp(rmargs[0], &rmargs[0]);
        exit(0);
    }

    if (pid1 > 0) {
        int r2;
        wait(&r2);
        return 1;
    }
    return 1;
}