#ifndef AUTHWINDOW_H
#define AUTHWINDOW_H

#include <QDialog>
#include <QDebug>
#include <QMessageBox>
#include <QFile>

#include <fstream>
#include <sstream>
#include <string>

#include <QNetworkAccessManager>
#include <QtNetwork/QNetworkReply>
#include <QtNetwork/QNetworkRequest>
#include <QtNetwork/QNetworkAccessManager>


#include <QJsonArray>
#include <QJsonValue>
#include <QJsonDocument>
#include <QJsonObject>

namespace Ui {
class AuthWindow;
}

class AuthWindow : public QDialog
{
    Q_OBJECT

public:
    explicit AuthWindow(QWidget *parent = nullptr);
    ~AuthWindow();

signals:
    void AuthorizedLogin(const QString& auth_login, const QString& auth_pass, const QString& auth_JWT);

private slots:
    void on_btnSignIn_clicked();

    void on_btnSignUp_clicked();

private:
    int Auth(const QString& login, const QString& pass);

    QNetworkAccessManager* manager;

    Ui::AuthWindow *ui;
    QString Login, Password, JWT;
    bool IsRemembered = false;
    QNetworkReply* MakeLoginRequest(const QString& login, const QString& pass);

    int ProcessLoginReply(QNetworkReply *, const QString&);

    const QString addressUpdate = "http://172.104.136.183/auth/update";
    const QString addressRegister = "http://172.104.136.183/auth/register";
    const QString addressRequest = "http://172.104.136.183/auth/request";
    const QString addressLogin = "http://172.104.136.183/auth/login";
};

#endif // AUTHWINDOW_H
