#include "../include/duload.h"

int download(const std::string& filename, const std::string& path, visFuncs* vis, int method) {
    std::vector<std::string> shards;
    std::string zipname = path + "/download_crowd";
    vis->Begin1(filename);
    int res = Merge(shards, zipname, path, vis);
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
int Merge(std::vector<std::string>& shards, const std::string& filename, const std::string& path, visFuncs* vis) {
    off_t curr_pt = 0;
    int current = 0, show = 0;
    int mode = 0666;
    (void)umask(0);
    int fdout;
    vis->SetField();

    //opening result file
    if ((fdout = open64 ((filename+".tar.gz").c_str(), O_RDWR | O_CREAT | O_TRUNC, mode )) < 0) {
            printf ("can't create %s for writing", filename.c_str());
            return 0;
        }
    for (int i = 0; i < shards.size(); i++) {
        struct stat statbuf;
        int fdin;
        //opening chunk file
        if ((fdin = open64 (shards[i].c_str(), O_RDONLY)) < 0) {
            printf("can't open %s for reading", shards[i].c_str());
            return 0;
        }
        //finding out size of chunk file
        if (fstat (fdin, &statbuf) < 0) {
            printf ("fstat error");
            return 0;
        }
        char *src, *dst;
        //mmapping chunk
        if ((src = (char*)mmap64 (0, statbuf.st_size, PROT_READ, MAP_SHARED, fdin, 0)) == (caddr_t) -1) {
            printf ("mmap error for input");
            return 0;
        }
        //making offset in result file
        if (lseek (fdout, curr_pt + statbuf.st_size - 1, SEEK_SET) == -1) {
            printf ("lseek error");
            return 0;
        }
        if (write (fdout, "", 1) != 1) {
            printf ("write error");
            return 0;
        }
        //mmapping result file
        if ((dst = (char*)mmap64 (0, statbuf.st_size, PROT_READ | PROT_WRITE, MAP_SHARED, fdout, curr_pt)) == (caddr_t) -1) {
            printf ("mmap error for output");
            return 0;
        }
        memcpy (dst, src, statbuf.st_size);
        munmap(src, statbuf.st_size);
        munmap(dst, statbuf.st_size);
        curr_pt += statbuf.st_size;
        close(fdin);

        current = (1.*i/shards.size()*100);
        vis->Next(show, current);
        show = current;

    }
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

