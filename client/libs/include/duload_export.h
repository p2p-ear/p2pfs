#pragma once

#include <string>
#include <vector>

typedef struct visFuncs {
    void (*SetField) ();
    void (*Next) (int, int);
    void (*Begin1) (const std::string&);
    void (*End1) (const std::string&);
    void (*Begin2) (const std::string&);
    void (*End2) (const std::string&);
} visFuncs;

extern int download(const std::string& filename, const std::string& path, visFuncs* vis, int method, std::string ip, unsigned long ringsz, std::string JWT, unsigned long nshards, std::string suff, unsigned long size);
extern int UploadFile(const std::string& filename, std::string name, const std::string& suff, bool remove, visFuncs* vis, unsigned long long shardLength, int method, std::string ip, unsigned long ringsz, std::string JWT);
extern int Merge(std::vector<std::string>& shards, const std::string& filename, const std::string& path, visFuncs* vis, std::string ip, unsigned long ringsz, std::string JWT, std::string suff, unsigned long size);
extern int unZIPFunc(const std::string& filename, const std::string& path, int method, visFuncs*);