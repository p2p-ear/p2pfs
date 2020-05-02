#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <fstream>
#include <QMessageBox>
#include <QButtonGroup>
#include <QDir>
#include <QListWidget>
#include <vector>

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

    void on_actionAuthorize_triggered();

    void on_actionQuit_triggered();

    void on_actionChange_user_triggered();

    void on_pushButton_clicked();

    void on_pushButton_2_clicked();

    void on_pushButton_3_clicked();

    void on_listWidget_itemDoubleClicked(QListWidgetItem *item);

    void on_listWidget_itemEntered(QListWidgetItem *item);

    void on_listWidget_itemActivated(QListWidgetItem *item);

private:
    Ui::MainWindow *ui;
    QString curr_path;
};
#endif // MAINWINDOW_H
