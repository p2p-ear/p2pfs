#ifndef DIALOG_H
#define DIALOG_H

#include <QDialog>

namespace Ui {
class Dialog;
}

class Dialog : public QDialog
{
    Q_OBJECT

public:
    explicit Dialog(QWidget *parent = nullptr);
    ~Dialog();

signals:
    void AuthorizedLogin(const QString& auth_login);

private slots:
    void on_signIn_clicked();

    void on_signUp_clicked();

private:
    Ui::Dialog *ui;
    QString Login;
    int Auth(const QString& login, const QString& pass);
    bool isRemembered = true;
};

#endif // DIALOG_H
