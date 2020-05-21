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

//powing
int pow(int bas, int value) {
    int res = 1;
    for (int i = 0; i < value; i++) {
        res *= bas;
    }
    return res;
}


//getting name by number of chunk
std::string getName(const std::string& name, int pos, int nSym) {
    std::string suff("_");
    int n = pow(23, nSym-1);
    for (int i = 0; i < nSym; i++) {
        //std::cout << "sym no" << i << " :"<<pos/n<<"\n";
        suff+=('a'+pos / n);
        pos %= n;
        n /= 23;
    }
    return name+suff;
}