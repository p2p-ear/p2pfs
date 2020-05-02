#include "../../libs/include/duload_export.h"

#include <filesystem>
#include <iostream>
#include <string>
#include <vector>
#include <unistd.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/msg.h>
#include <string.h>
#include <sys/sem.h>
#include <sstream>
#include <signal.h>

namespace fs = std::filesystem;

int semnum = 0;


void handler(int signum) {
    semctl(semnum, IPC_RMID, 0);
    exit(0);
}

void EndShard(const std::string& res) {
    std::cout << "  ";
    for (int i = 0; i < 101; i++) {
        std::cout << "\b";
    }
    for (int i = 0; i < 101; i++) {
        std::cout << " ";
    }
    for (int i = 0; i < 101; i++) {
        std::cout << "\b";
    }
}

void Next(int num, int curr) {
    //std::cout << curr << "\n";
    int to_add = curr - num;
    for (unsigned long long i = 0; i < to_add; i++) {
            std::cout << "#";
            std::cout.flush();
    }
}

void SetField() {
    /*std::cout << "[";
    for (int i = 0; i < 100; i++) {
        std::cout << " ";
    }
    std::cout<<"]";
    for (int i = 0; i < 101; i++) {
        std::cout << "\b";
    }*/
}

void StartZip(const std::string& filename) {
    //std::cout << filename << " is being ZIPed\n";
    std::cout << "[";
    for (int i = 0; i < 100; i++) {
        std::cout << ".";
    }
    std::cout<<"]";
    for (int i = 0; i < 101; i++) {
        std::cout << "\b";
    }
    std::cout.flush();
}

void EndZip(const std::string& filename) {
    //std::cout << filename << " is ZIPed\n";
}

unsigned long long EvaluateSize(std::vector<fs::path>& args, const std::string& start_path) {
    unsigned long long res = 0;
    for (auto& arg : args) {
        if (fs::exists(arg)) {
            if (fs::is_directory(arg)) {
                for (const auto &entry_point : fs::recursive_directory_iterator(arg, fs::directory_options::skip_permission_denied)) {
                    if (fs::is_regular_file(entry_point.path())) {
                        res += entry_point.file_size();
                    }
                }
            } else {
                res += fs::file_size(arg);
            }
        } else {
            std::cout << "No such file or directory :\"" << arg.string() << "\"\n";
        }
    }
    return res;
}

int checkAuth() {
    return 1;
}

int askAuth() {
    return 1;
}

std::string GetIp() {
    return "192.168.0.0:9000";
}

unsigned long GetRingSz() {
    return 1000;
}

int askForLoad(unsigned long long size) {
    if (size <= 10000000000) {
        return 1;
    } else {
        return 0;
    }
}

int checkSize(unsigned long long size) {
    return askForLoad(size);
}

int main(int argc, char** argv) {
    //initializing application
    signal(SIGKILL, handler);
    signal(SIGINT, handler);
    signal(SIGTERM, handler);

    int method = 1;
    bool remove = false;

    if ((semnum = semget(123456, 1, IPC_CREAT | IPC_EXCL)) == -1) {
        printf("Programm is already launched\n");
        return 1;
    }

    //args parsing
    if (argc > 1) {
        for (int i = 1; i < argc; i++) {
            if (strlen(argv[i]) > 1 && argv[i][0] == '-') {
                if (argv[i][1] == 'm') {
                    remove = true;
                } else if (argv[i][1] == 'c') {
                    if (strlen(argv[i]) > 2) {
                        std::stringstream mode(&argv[i][2]);
                        mode >> method;
                        //std::cout << method << "\n";
                    } 
                } else {
                std::cout << "Unknown parameter \"" << argv[i] <<"\"\n";
                semctl(semnum, 0, IPC_RMID);
                return -1;
                }
            } else {
                std::cout << "Unknown parameter \"" << argv[i] <<"\"\n";
                semctl(semnum, 0, IPC_RMID);
                return -1;
            }
        }
        
    }
    //asking for auth

    if (checkAuth() != 1 && askAuth() != 1) {
        semctl(semnum, 0, IPC_RMID);
        return -1;
    } 


    //choosing files to load
    std::cout << "Choose files to load:\n";

    std::string files, file;
    getline(std::cin, files);
    std::stringstream filestream(files);

    std::vector<fs::path> args;

    while(filestream >> file) {
        if (fs::exists(fs::path(file))) {
            args.push_back(fs::path(file));
        }
    }    
    char buf[4096];
    getcwd(buf, 4096);
    std::string start_path(buf);
    unsigned long long size = EvaluateSize(args, start_path);
    std::cout << size << " bytes\n";

    //loading
    if (size > 0) {
        if (checkSize(size) == 1) {
            std::cout << "Are you sure to load files:\n";
            for (const auto& item : args) {
                std::cout << item << "\n";
            }
            std::cout << "?[y/N]\n";
            std::string answer;
            std::cin >> answer;
            if (answer == "Y" || answer == "y") {

                visFuncs vis;
                vis.SetField = SetField;
                vis.Next = Next;
                vis.End1 = EndShard;
                vis.Begin2 = StartZip;
                vis.End2 = EndZip;

                for (const auto& item : args) {
                    std::cout << item << " :\n";
                    unsigned long RingSz = GetRingSz();
                    std::string Ip = GetIp();
                    int result = UploadFile(item, "", remove, &vis, 1600, method, Ip, RingSz);
                    if (result == 1) {
                        std::cout << "is loaded]\n";
                    } else {
                        std::cout << "is not loaded]\n";
                    }
                }
                semctl(semnum, 0, IPC_RMID);
                return 1;
            } else {
                std::cout << "Abort loading\n";
                semctl(semnum, 0, IPC_RMID);
                return -1;
            }
        } else {
            std::cout << "Cannot load files of this size\n";
            semctl(semnum, 0, IPC_RMID);
            return -1;
        }
    } else {
        std::cout << "Total size of files is 0 bytes\n";
    }

    semctl(semnum, 0, IPC_RMID);
}