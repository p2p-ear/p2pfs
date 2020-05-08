#include "dialog.h"
#include "ui_dialog.h"
#include <QDebug>
#include <QMessageBox>

#include <fstream>
#include <sstream>
#include <string>

Dialog::Dialog(QWidget *parent) :
    QDialog(parent),
    ui(new Ui::Dialog)
{
    ui->setupUi(this);
    std::fstream authConfig(".auth");
    std::string rememberStr;
    getline(authConfig, rememberStr);
    std::stringstream ss(rememberStr);
    ss >> isRemembered;
    if (isRemembered) {
        std::string tmp;
        getline(authConfig, tmp);
        Login = QString(tmp.c_str());
        ui->LineLogin->setText(Login);
        ui->checkBox->setCheckState(Qt::CheckState::Checked);
    }
    authConfig.close();
}

Dialog::~Dialog()
{
    delete ui;
}

void Dialog::on_signIn_clicked()
{
    qDebug() << "sign_in\n";
    QString login = ui->LineLogin->text(), pass = ui->LinePass->text();
    if (Auth(login, pass) == 1) {
        std::fstream authConfig(".auth");
        authConfig.clear();
        if (ui->checkBox->checkState() == Qt::CheckState::Checked) {
            isRemembered = true;
            authConfig << isRemembered << "\n";
            Login = ui->LineLogin->text();
            authConfig << Login.toStdString();
        } else {
            isRemembered = false;
            authConfig << isRemembered << "\n";
        }
        authConfig.close();
        emit AuthorizedLogin(Login);
        QMessageBox::information(this, "Authoriztion", "Authorization successful");
        this->close();
    } else {
        QMessageBox::information(this, "Authoriztion", "Authorization failed");
    }
}

void Dialog::on_signUp_clicked()
{
    qDebug() << "sign_up\n";
}

int Dialog::Auth(const QString &login, const QString &pass)
{
    if (login == "sanya" && pass == "123") {
        return 1;
    } else {
        return 0;
    }
}
