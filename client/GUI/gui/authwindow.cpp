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

    manager = new QNetworkAccessManager;
}

AuthWindow::~AuthWindow()
{
    delete ui;
}

void AuthWindow::on_btnSignIn_clicked() {
    qDebug() << "sign_in\n";
    QString login = ui->lineLogin->text(), pass = ui->linePasswd->text();
    if (Auth(login, pass) == 1) {
        QFile f(".auth");
        f.remove();
        f.open(QIODevice::ReadWrite);
        f.close();
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
        emit AuthorizedLogin(login, pass, JWT);
        QMessageBox::information(this, "Authoriztion", "Authorization successful");
        this->close();
    } else {
        QFile f(".auth");
        f.remove();
        f.open(QIODevice::ReadWrite);
        f.close();
        std::fstream authConfig(".auth");

        authConfig.clear();
        IsRemembered = false;
        authConfig << IsRemembered << "\n";
        authConfig.close();

        ui->checkBox->setCheckState(Qt::CheckState::Unchecked);
    }
}

void AuthWindow::on_btnSignUp_clicked() {

}

QNetworkReply* AuthWindow::MakeLoginRequest(const QString& login, const QString& pass) {
    QJsonObject jObj;
    QJsonObject body;
    jObj.insert("email", login);
    jObj.insert("password", pass);
    jObj.insert("body", body);
    qDebug() << jObj.keys().size();
    QJsonDocument jDoc(jObj);
    QNetworkRequest req;
    req.setUrl(QUrl(addressLogin));
    qDebug() << jDoc.toJson();
    req.setRawHeader("Content-Type","application/json");
    auto reply = manager->post(req, jDoc.toJson());
    QEventLoop loop;
    connect(reply, SIGNAL(finished()), &loop, SLOT(quit()));
    loop.exec();
    return reply;
}

int AuthWindow::ProcessLoginReply(QNetworkReply * reply, const QString& login) {
    QByteArray responseData = reply->readAll();
    QJsonDocument doc(QJsonDocument::fromJson(responseData));
    QJsonObject rep = doc.object();
    if (rep["status"].toInt() != 0) {//add int value for error processing
        QMessageBox::warning(this, "Error", rep["message"].toString());
        return 1;
    } else if (rep["email"] != login) { //add logout!!!!!!!!!!!
        qDebug() << login;
        qDebug() << rep["email"];
        QMessageBox::warning(this, "Authoriztion", "Authorization failed");
    } else {
        if (rep["status"].toInt() == 0) { //authorization success
            JWT = rep["auth_token"].toString();
            return 0;
        }
    }
    return 0;
}


int AuthWindow::Auth(const QString &login, const QString &pass) {
    QNetworkReply * reply = MakeLoginRequest(login, pass);
    if (ProcessLoginReply(reply, login) == 0) {
        return 1;
    } else {
        return 0;
    }
}
