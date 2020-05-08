#include <sys/sem.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <signal.h>
#include <string.h>
#include <iostream>
#include <sstream>
#include <iostream>


#include "../../libs/include/duload_export.h"

void SetField () {
    return ;
}
void Next (int num, int curr) {
    int to_add = curr - num;
    for (unsigned long long i = 0; i < to_add; i++) {
        std::cout << "#";
        std::cout.flush();
    }
}
void Begin1 (const std::string&) {
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
void End1 (const std::string& filename) {
    return ;
}
void Begin2 (const std::string&) {
    return;
}
void End2 (const std::string&) {
    return;
}

int semnum = 0;


void handler(int signum) {
    semctl(semnum, IPC_RMID, 0);
    exit(0);
}

int checkAuth() {
    return 1;
}

int askAuth() {
    return 1;
}

int main(int argc, char** argv) {
    //initializing application
    signal(SIGKILL, handler);
    signal(SIGINT, handler);
    signal(SIGTERM, handler);
    int method = 1;



    if ((semnum = semget(123456, 1, IPC_CREAT | IPC_EXCL)) == -1) {
        printf("Programm is already launched\n");
        return 1;
    }

    //args parsing
    /*if (argc > 1) {
        for (int i = 1; i < argc; i++) {
            if (strlen(argv[i]) > 1 && argv[i][0] == '-') {
                if (argv[i][1] == 'c') {
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
        
    }*/


    //check auth
    if (checkAuth() != 1 && askAuth() != 1) {
        semctl(semnum, 0, IPC_RMID);
        return -1;
    } 

    std::vector<std::string> shards;
    for (int i = 1; i < argc; i++) {
        shards.push_back(std::string(argv[i]));
    }

    visFuncs vis;
    vis.Begin1 = Begin1;
    vis.Begin2 = Begin2;
    vis.End1 = End1;
    vis.End2 = End2;
    vis.Next = Next;
    vis.SetField = SetField;

    std::string path;
    std::cout << "Select path for downloading:\n";
    std::cin >> path;
    std::string filename = path + "/download_crowd";
    int res = Merge(shards, filename, path, &vis);
    
    filename += ".tar.gz";
    int result = unZIPFunc(filename, path, method, &vis);
    //RemoveFile(filename);

    semctl(semnum, 0, IPC_RMID);
}