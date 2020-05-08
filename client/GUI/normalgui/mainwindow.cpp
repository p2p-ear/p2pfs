#include "mainwindow.h"
#include "ui_mainwindow.h"
#include "dialog.h"



MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    ui->actionUsername_something_com->setDisabled(true);
    ui->actionQuit->setDisabled(true);
    ui->actionChange_user->setDisabled(true);
    ui->listWidget->setSelectionMode(QAbstractItemView::MultiSelection);
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::setAuthLogin(const QString &auth_login) {
    ui->actionQuit->setEnabled(true);
    ui->actionChange_user->setEnabled(true);
    ui->actionUsername_something_com->setEnabled(true);
    ui->actionUsername_something_com->setText(auth_login);
}


void MainWindow::on_actionAuthorize_triggered() {
    Dialog *auth = new Dialog(this);
    connect(auth, SIGNAL(AuthorizedLogin(QString)), this, SLOT(setAuthLogin(QString)));
    auth->setModal(true);
    auth->exec();
    disconnect(auth, SIGNAL(AuthorizedLogin(QString)), this, SLOT(setAuthLogin(QString)));
    delete auth;
}

void MainWindow::on_actionQuit_triggered() {
    QMessageBox::StandardButton btn =  QMessageBox::question(this, "Confirm action", "Are you sure to logout?", QMessageBox::Yes | QMessageBox::No);
    if (btn == QMessageBox::Yes) {
        ui->actionChange_user->setDisabled(true);
        ui->actionQuit->setDisabled(true);
        ui->actionUsername_something_com->setText("username@someth.ing");
        ui->actionUsername_something_com->setDisabled(true);
        std::fstream authFile(".auth");
        authFile.clear();
        authFile << false << "\n";
        authFile.close();
    }
}

void MainWindow::on_actionChange_user_triggered() {
    QMessageBox::StandardButton btn =  QMessageBox::question(this, "Confirm action", "Are you sure to change user?", QMessageBox::Yes | QMessageBox::No);
    if (btn == QMessageBox::Yes) {
        ui->actionChange_user->setDisabled(true);
        ui->actionQuit->setDisabled(true);
        ui->actionUsername_something_com->setText("username@someth.ing");
        ui->actionUsername_something_com->setDisabled(true);
        std::fstream authFile(".auth");
        authFile.clear();
        authFile << false << "\n";
        authFile.close();

        on_actionAuthorize_triggered();
    }
}

void MainWindow::on_pushButton_clicked() {
    QString path = ui->lineEdit->text();
    curr_path = path;
    QDir dir(path);
    ui->listWidget->clear();
    for (const auto& item : dir.entryInfoList()) {
        ui->listWidget->addItem(item.fileName());
    }
}

void MainWindow::on_pushButton_2_clicked() {
    QListWidgetItem *item = ui->listWidget->currentItem();
    QString dst = item->text();
    QDir curr_dir(curr_path);
    curr_dir.cd(dst);
    curr_path+="/"+dst;
    ui->listWidget->clear();
    for (const auto& item : curr_dir.entryInfoList()) {
        ui->listWidget->addItem(item.fileName());
    }
}

void MainWindow::on_pushButton_3_clicked()
{
    ui->textBrowser->clear();
    QString res;
    for (const auto& item : ui->listWidget->selectedItems()) {
        res += item->text() + "\n";

    }
    ui->textBrowser->setText(res);
}

void MainWindow::on_listWidget_itemDoubleClicked(QListWidgetItem *item) {
    QString dst = item->text();
    QDir curr_dir(curr_path);
    curr_dir.cd(dst);
    curr_path+="/"+dst;
    ui->listWidget->clear();
    for (const auto& item : curr_dir.entryInfoList()) {
        ui->listWidget->addItem(item.fileName());
    }
}

//wrong
void MainWindow::on_listWidget_itemEntered(QListWidgetItem *item) {
    QString dst = item->text();
    QDir curr_dir(curr_path);
    curr_dir.cd(dst);
    curr_path+="/"+dst;
    ui->listWidget->clear();
    for (const auto& item : curr_dir.entryInfoList()) {
        ui->listWidget->addItem(item.fileName());
    }
}

//enter moving
void MainWindow::on_listWidget_itemActivated(QListWidgetItem *item)
{
    QString dst = item->text();
    QDir curr_dir(curr_path);
    curr_dir.cd(dst);
    curr_path+="/"+dst;
    ui->listWidget->clear();
    for (const auto& item : curr_dir.entryInfoList()) {
        ui->listWidget->addItem(item.fileName());
    }
}
