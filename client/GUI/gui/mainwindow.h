#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QListWidget>
#include <QString>
#include <QDir>
#include <QMessageBox>
#include <stack>
#include <set>
#include <filesystem>
#include <QTableWidget>

#include <QNetworkAccessManager>
#include <QtNetwork/QNetworkReply>
#include <QtNetwork/QNetworkRequest>
#include <QtNetwork/QNetworkAccessManager>

#include "json.h"

#include "../../libs/include/duload_export.h"

namespace fs = std::filesystem;


QT_BEGIN_NAMESPACE
namespace Ui { class MainWindow; }
QT_END_NAMESPACE

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

private slots:
    void setAuthLogin(const QString& auth_login, const QString& auth_pass, const QString& auth_JWT);

    void on_btnBack_clicked();

    void on_btnForward_clicked();

    void on_btnUp_clicked();

    void on_btnHome_clicked();

    void on_btnAdd_clicked();

    void on_btnRemove_clicked();

    void on_btmClear_clicked();

    void on_listWidget_itemDoubleClicked(QListWidgetItem *item);

    void on_listWidget_itemActivated(QListWidgetItem *item);

    void on_btnPath_clicked();

    void on_btnUpload_clicked();

    void on_actionLogout_triggered();

    void on_actionChange_User_triggered();

    void on_actionAuthorize_triggered();

    void on_actionUser_Options_triggered();

    void on_actionUsername_some_thing_triggered();

    void on_btnCd_clicked();

    void updateTable(const QString&);

    void updateTable2(const QString&);

    void updateTable3(const std::vector<MDfile>&);

    void updateTable4(const MDfile&);

    void on_tableWidget_itemActivated(QTableWidgetItem *item);

    void on_btnCd_2_clicked();

    void on_btnAdd2_clicked();

    void on_btnRemove2_clicked();

    void on_btmClear2_clicked();

    void on_pushButton_clicked();

    void on_pushButton_2_clicked();

    void on_btnAddCoins_clicked();

    void on_btnUodateJson_clicked();

    void on_btnGetCoins_clicked();

private:
    const QString addressUpdate = "http://172.104.136.183/auth/update";
    const QString addressRegister = "http://172.104.136.183/auth/register";
    const QString addressRequest = "http://172.104.136.183/auth/request";



    Ui::MainWindow *ui;
    QString current_path;
    std::stack<QString> uploadBack, uploadForward;
    std::set<QString> uploadset, downloadset;
    unsigned long long totalSize = 0;
    bool is_authorised = false;

    MyDiskFs FS;

    QString Login, Password, JWT;

    QNetworkAccessManager * manager;
    int MakeReqRequest(QJsonObject& body, int type);

    unsigned long long EvaluateSize(std::vector<fs::path>& args, const std::string& start_path);
};
#endif // MAINWINDOW_H
