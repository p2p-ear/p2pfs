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
    void setAuthLogin(const QString& auth_login);

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

private:
    Ui::MainWindow *ui;
    QString current_path;
    std::stack<QString> uploadBack, uploadForward;
    std::set<QString> uploadset;
    unsigned long long totalSize = 0;

    unsigned long long EvaluateSize(std::vector<fs::path>& args, const std::string& start_path);
};
#endif // MAINWINDOW_H
