#include "../include/duload.h"


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


//getting num of chunks by size of file
int getNumPieces(unsigned long long size, int mode) {
    int res = size % mode;
    return size / mode + (res > 0 ? 1 : 0);
}


//deviding into chunks
int shardFile(const std::string &filename, visFuncs* vis, unsigned long long shardLength, std::string ip, unsigned long ringsz, std::string JWT) {
    vis->SetField();
    int show = 0, current = 0;

    umask(0);
    unsigned long long start_pt = 0;
    int mode = 0x0777;
    int fdin, fdout;
    unsigned long long shardSize = shardLength*PG_SIZE;


    //opening the file
    if ((fdin = open64 (filename.c_str(), O_RDONLY)) < 0) {
        printf("can't open %s for reading", filename.c_str());
        return 0;
    }

    //printf("1\n");


    //checking size of the file
    struct stat statbuf;

    if (fstat (fdin, &statbuf) < 0) {
        printf ("fstat error");
        return 0;
    }

    //printf("2\n");


    //getting num of chunks
    int num = getNumPieces(statbuf.st_size, shardSize);
    std::string name(filename);
    off_t curr_pt = start_pt;

    int nSym = ceil(log(num)/log(23));

    //printf("3\n");

    GoString IP = {ip.c_str(), ip.length()};
    GoString wJWT = {JWT.c_str(), JWT.length()};

    
    for (int i = 0; i < num; i++) {
        unsigned long long curr_size = statbuf.st_size - curr_pt < shardSize ? statbuf.st_size - curr_pt : shardSize;
        
        //printf("4\n");
        std::string shardName = getName(name, i, nSym);
        //printf("5\n");

        GoString fname = {shardName.c_str(), shardName.length()};


        char * src;
        //printf("8\n");

        //mmapping file
        if ((src = (char*)mmap64 (0, curr_size, PROT_READ, MAP_SHARED, fdin, curr_pt)) == (caddr_t) -1) {
            printf ("mmap error for input");
            return 0;
        }
        //printf("9\n");

        GoSlice fcontent_clice = {
            data: src,
            len: curr_size,
            cap: curr_size
        };

        UploadFileRSC(IP, fname, ringsz, fcontent_clice, wJWT);

        munmap(src, curr_size);

        //printf("12\n");

        curr_pt += shardSize;
        
        current = (1.*i/num*100);
        vis->Next(show, current);
        show = current;
        
    }

    vis->End1(filename);
    return 1;
}

std::string TarGz(const fs::path filename) {
    std::vector<char*> params;
    std::string res;
    res = filename.string()+".tar.gz";
    std::string command("tar");
    std::string flags("-czf");
    std::string flags2("-C");
    std::string flags3("--absolute-names");
    std::string path = filename.parent_path().string();

    params.push_back((char*)command.c_str());
    //params.push_back((char*)flags3.c_str());
    params.push_back((char*)flags2.c_str());
    params.push_back((char*)path.c_str());
    params.push_back((char*)flags.c_str());
    params.push_back((char*)res.c_str());
    params.push_back((char*)filename.filename().c_str());
    params.push_back(nullptr);

    int pid = fork();
    if (pid == -1) {
        return "";
    }
    if (pid == 0) { //child
        execvp(params[0], &params[0]);
        exit(0);
    }
    if (pid > 0) {
        int r;
        wait(&r);
    }
    return res;
}

std::string ZIPFunc(const fs::path& filename, int method) {
    if (method == 1) {
        return TarGz(filename);
    } else {
        return "";
    }
}



std::string Pbzip2(const std::string& filename) {
    std::vector<char*> params;
    std::string res;
    res = filename+".bz2";
    std::string command("pbzip2");
    std::string flags("-k");
    std::string flags2("-p4");

    params.push_back((char*)command.c_str());
    params.push_back((char*)flags.c_str());
    params.push_back((char*)flags2.c_str());
    params.push_back((char*)filename.c_str());
    params.push_back(nullptr);

    int pid = fork();
    if (pid == -1) {
        return "";
    }
    if (pid == 0) { //child
        execvp(params[0], &params[0]);
        exit(0);
    }
    if (pid > 0) {
        int r;
        wait(&r);
    }
    return res;
}

void RollDir(const std::string&) {
    return;
}

int UploadFile(const std::string& filename, const std::string& suff, bool remove, visFuncs* vis, unsigned long long shardLength, int method, std::string ip, unsigned long ringsz, std::string JWT) { //zero if fail, else 1
    bool isDir = false;
    vis->Begin2(filename);
    fs::path file(filename);
    std::string ZipedFile;
    //checking if file exists
    if (!fs::exists(file)) {
        return 0;
    }
    //checking if file is a directory
    if (fs::is_directory(file)) {
        isDir = true;
    } 

    //zipping
    ZipedFile = ZIPFunc(file, 1);
    if (ZipedFile == "") {
        return 0;
    }
    vis->End2(filename);

    //chunking
    int result = shardFile(ZipedFile, vis, shardLength, ip, ringsz, JWT);
    //std::cout << ZipedFile << "\n";

    //removing temp files
    RemoveFile(ZipedFile);
    //removing the file if is requeired
    if (result == 1 && remove == true) {
        std::cout << "removing\n";
        RemoveFile(file.string());
    }
    return result;
}