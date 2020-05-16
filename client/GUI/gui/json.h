#pragma once

#include <QJsonArray>
#include <QJsonValue>
#include <QJsonDocument>
#include <QJsonObject>
#include <QFile>
#include <vector>
#include <QDebug>

#define T_NAME "Name"
#define T_ISDIR "IsDir"
#define T_SIZE "Size"
#define T_CHILD "Child"

struct MDfile {
    QString Name;
    bool isDir;
    unsigned long long Size;
};

class MyDiskFs {
public:
    int Load() {
        QFile jFile("test.json");
        if (jFile.open(QIODevice::ReadOnly | QFile::Text)) {
            qDebug() << jFile.size();
            QJsonDocument doc = QJsonDocument::fromJson(QByteArray(jFile.readAll()));
            root = QJsonObject(doc.object());
            jFile.close();
            return 1;
        } else {
            return 0;
        }
    }

    int Load(QJsonObject res) {
        root = res;
        return 1;
    }

    int Cd(const QString& aim) {
        QJsonObject currentDir = GetObj(aim);
        if (currentDir[T_NAME].toString() == "") {
            return 0;
        }
        path = aim;
        return 1;
    }

    std::vector<MDfile> Ls() {
        std::vector<MDfile> res;
        QJsonObject subdir = GetObj(path)[T_CHILD].toObject();
        for (const auto& dir : subdir.keys()) {
            QJsonObject obj = subdir[dir].toObject();
            MDfile fl = {obj[T_NAME].toString(), obj[T_ISDIR].toBool(), (unsigned long long)obj[T_SIZE].toInt()};
            res.push_back(fl);
        }
        return res;
    }

    QString GetCurrPath() {
        return path;
    }

private:
    QJsonObject root;
    QString path = "/";
    QJsonObject ret;


    QJsonObject GetObj(const QString& aim) {
        QString pth = aim.toStdString().substr(1).c_str();
        QList list = pth.split('/');
        QJsonObject ptr = root;
        for (const auto& dir : list) {
            if (dir == "") {
                continue;
            }
            if (ptr[T_ISDIR].toBool() == true) {
                QJsonObject subdir = ptr[T_CHILD].toObject();
                if (subdir.contains(dir) && subdir[dir].toObject()[T_ISDIR].toBool() == true) {
                    ptr = subdir[dir].toObject();
                } else {
                    ret.insert(T_NAME, "");
                    return ret;
                }
            } else {
                ret.insert(T_NAME, "");
                return ret;
            }
        }
        return ptr;
    }
};
