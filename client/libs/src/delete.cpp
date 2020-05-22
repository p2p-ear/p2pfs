#include "../include/duload.h"

int Delete(std::string ip, std::string filename, unsigned long ringsz, std::string JWT, unsigned long nshards) {

    GoString IP = {ip.c_str(), ip.length()};
    GoString dJWT = {JWT.c_str(), JWT.length()};

    GoString codename = {(filename+"_KEY").c_str(), (filename+"_KEY").length()};
    //deleting codes
    DeleteFileRSC(IP, codename, ringsz, dJWT);
    std::vector<std::string> shards;
    GetNames(filename, shards, nshards);

    for (unsigned long i = 0; i < nshards; i++) {
        std::string currName = shards[i];

        GoString Name = {currName.c_str(), currName.length()};
        DeleteFileRSC(IP, Name, ringsz, dJWT);

    }
}