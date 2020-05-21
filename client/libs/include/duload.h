#pragma once

#include "libc_interface.h"
#include "cypher.h"

#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <vector>
#include <string>
#include <iostream>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <math.h>
#include <sys/wait.h>
#include <filesystem>

#define PG_SIZE 4096llu

typedef struct visFuncs {
    void (*SetField) ();
    void (*Next) (int, int);
    void (*Begin1) (const std::string&);
    void (*End1) (const std::string&);
    void (*Begin2) (const std::string&);
    void (*End2) (const std::string&);
} visFuncs;

namespace fs = std::filesystem;


//upload
int pow(int bas, int value);
std::string getName(const std::string& name, int pos, int nSym);
int getNumPieces(unsigned long long size, int mode);
int shardFile(const std::string &filename, std::string name, visFuncs* vis, unsigned long long shardLength, std::string ip, unsigned long ringsz, std::string JWT);
int UploadFile(const std::string& filename, std::string name, const std::string& suff, bool remove, visFuncs* vis, unsigned long long shardLength, int method, std::string ip, unsigned long ringsz, std::string JWT);
std::string ZIPFunc(const fs::path& filename, int method);

//zip funcs
std::string TarGz(const fs::path filename);
std::string Pbzip2(const std::string& filename);


//download
int download(const std::string& filename, const std::string& path, visFuncs* vis, int method, std::string ip, unsigned long ringsz, std::string JWT, unsigned long nshards, std::string suff, unsigned long size);
int Merge(std::vector<std::string>& shards, const std::string& filename, const std::string& path, visFuncs* vis, std::string ip, unsigned long ringsz, std::string JWT, std::string suff, unsigned long size);
int unZIPFunc(const std::string& filename, const std::string& path, int method, visFuncs*);

//unzip funcs
int unTarGz(const std::string& filename, const std::string& path);


//common
int LastSlash(const std::string& path);
int RemoveFile(const std::string& filename);


