#include "../include/duload.h"

void GetNames(const std::string& filename, std::vector<std::string>& shards, unsigned long nshards) {
    int nSym = ceil(log(nshards)/log(23));

    for (unsigned long i = 0; i < nshards; i++) {
        shards.push_back(getName(filename, i, nSym));
    }
}

int download(const std::string& filename, const std::string& path, visFuncs* vis, int method, std::string ip, unsigned long ringsz, std::string JWT, unsigned long nshards, std::string suff, unsigned long size) {
    std::vector<std::string> shards;
    GetNames(filename, shards, nshards);
    std::string zipname = path + "/download_crowd";
    vis->Begin1(filename);
    int res = Merge(shards, zipname, path, vis, ip, ringsz, JWT, suff, size);
    vis->End1(filename);
    if (res > 0) {
        zipname += ".tar.gz";
        vis->Begin2(filename);
        int result = unZIPFunc(zipname, path, method, vis);
        vis->End2(filename);
        //RemoveFile(filename);
        return result;
    } else {
        return 0;
    }
}

//merging chunks into file
int Merge(std::vector<std::string>& shards, const std::string& filename, const std::string& path, visFuncs* vis, std::string ip, unsigned long ringsz, std::string JWT, std::string suff, unsigned long size) {
    off_t curr_pt = 0;
    int current = 0, show = 0;
    int mode = 0666;
    (void)umask(0);
    int fdout;
    vis->SetField();

    GoSlice buff = {
        //Allocate at least 9 more bytes than nescessary
        data: (GoInt8*) malloc(size*sizeof(GoInt8)),
        len: size,
        cap: size
    };

    GoString rJWT = {JWT.c_str(), JWT.length()};
    GoString IP = {ip.c_str(), ip.length()};


    GoString codename = {(filename+"_KEY").c_str(), (filename+"_KEY").length()};

    //getting code
    DownloadFileRSC(IP, codename, ringsz, buff, rJWT);
    XORcypher decoder(size, (char*)buff.data);

    //opening result file
    if ((fdout = open64 ((filename+".tar.gz").c_str(), O_RDWR | O_CREAT | O_TRUNC, mode )) < 0) {
            printf ("can't create %s for writing", filename.c_str());
            return 0;
        }
    for (int i = 0; i < shards.size(); i++) {
        int fdin;
        
        GoString fname = {(suff+shards[i]).c_str(), (suff+shards[i]).length()};
        

        char *dst;
        //mmapping chunk

        GoInt remainingSpace = DownloadFileRSC(IP, fname, ringsz, buff, rJWT);

        //making offset in result file
        if (lseek (fdout, curr_pt + (size-remainingSpace) - 1, SEEK_SET) == -1) {
            printf ("lseek error");
            return 0;
        }
        if (write (fdout, "", 1) != 1) {
            printf ("write error");
            return 0;
        }
        //mmapping result file
        if ((dst = (char*)mmap64 (0, (size-remainingSpace), PROT_READ | PROT_WRITE, MAP_SHARED, fdout, curr_pt)) == (caddr_t) -1) {
            printf ("mmap error for output");
            return 0;
        }

        decoder((char*)buff.data);

        memcpy (dst, buff.data, (size-remainingSpace));
        munmap(dst, (size-remainingSpace));
        curr_pt += (size-remainingSpace);

        current = (1.*i/shards.size()*100);
        vis->Next(show, current);
        show = current;

    }

    //decode
    close(fdout);
    return 1;
}

int unTarGz(const std::string& filename, const std::string& path) {
    std::vector<char*> params;
    std::string res = filename.substr(0, filename.length() - 7);
//    std::string path = res.substr(0, LastSlash(res));
    std::string command("tar");
    std::string flags("-xf");
    std::string flags2("-C");

    params.push_back((char*)command.c_str());
    params.push_back((char*)flags2.c_str());
    params.push_back((char*)path.c_str());
    params.push_back((char*)flags.c_str());
    params.push_back((char*)filename.c_str());
    params.push_back(nullptr);

    int pid1 = fork();
    if (pid1 == -1) {
        return 0;
    }
    if (pid1 == 0) { //child
        execvp(params[0], &params[0]);
        exit(0);
    }
    if (pid1 > 0) {
        int r1;
        wait(&r1);
    }
    return 1;
}

int unZIPFunc(const std::string& filename, const std::string& path, int method, visFuncs* vis) {
    int fdin;
    if (!fs::exists(fs::path(filename))) {
        return 0;
    }

    if (unTarGz(filename, path) == 0) {
        return 0;
    }

    RemoveFile(filename);
    return 1;
}

