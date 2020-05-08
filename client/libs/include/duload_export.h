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

extern int download(const std::string& filename, const std::string& path, visFuncs*, int method);
extern int UploadFile(const std::string& filename, const std::string& suff, bool remove, visFuncs* vis, unsigned long long shardLength, int method, std::string ip, unsigned long ringsz);
extern int Merge(std::vector<std::string>& shards, const std::string& filename, const std::string& path, visFuncs*);
extern int unZIPFunc(const std::string& filename, const std::string& path, int method, visFuncs*);