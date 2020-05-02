#include "authwindow.h"
#include "ui_authwindow.h"

AuthWindow::AuthWindow(QWidget *parent) :
    QDialog(parent),
    ui(new Ui::AuthWindow)
{
    ui->setupUi(this);
    std::fstream authConfig(".auth");
    std::string rememberStr;
    getline(authConfig, rememberStr);
    std::stringstream ss(rememberStr);
    ss >> IsRemembered;
    if (IsRemembered) {
        std::string tmp;
        getline(authConfig, tmp);
        Login = QString(tmp.c_str());
        ui->lineLogin->setText(Login);
        ui->checkBox->setCheckState(Qt::CheckState::Checked);
    }
    authConfig.close();
}

AuthWindow::~AuthWindow()
{
    delete ui;
}

void AuthWindow::on_btnSignIn_clicked() {
    qDebug() << "sign_in\n";
    QString login = ui->lineLogin->text(), pass = ui->linePasswd->text();
    if (Auth(login, pass) == 1) {
        std::fstream authConfig(".auth");
        authConfig.clear();
        if (ui->checkBox->checkState() == Qt::CheckState::Checked) {
            IsRemembered = true;
            authConfig << IsRemembered << "\n";
            Login = ui->lineLogin->text();
            authConfig << Login.toStdString();
        } else {
            IsRemembered = false;
            authConfig << IsRemembered << "\n";
        }
        authConfig.close();
        emit AuthorizedLogin(Login);
        QMessageBox::information(this, "Authoriztion", "Authorization successful");
        this->close();
    } else {
        QMessageBox::information(this, "Authoriztion", "Authorization failed");
    }
}

void AuthWindow::on_btnSignUp_clicked() {

}

int AuthWindow::Auth(const QString &login, const QString &pass) {
    if (login == "sanya" && pass == "123") {
        return 1;
    } else {
        return 0;
    }
}
