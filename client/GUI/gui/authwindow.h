#ifndef AUTHWINDOW_H
#define AUTHWINDOW_H

#include <QDialog>
#include <QDebug>
#include <QMessageBox>

#include <fstream>
#include <sstream>
#include <string>

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


    Ui::AuthWindow *ui;
    QString Login, Password;
    bool IsRemembered = false;
};

#endif // AUTHWINDOW_H
