#include "mainwindow.h"
#include "ui_mainwindow.h"
#include "authwindow.h"

void vis1() {
    return;
}
void vis2(const std::string& res) {
    return;
}

void vis3(int a, int b) {
    return;
}

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    ui->actionUsername_some_thing->setDisabled(true);
    ui->actionLogout->setDisabled(true);
    ui->actionChange_User->setDisabled(true);
    ui->actionUser_Options->setDisabled(true);
    ui->listWidget->setSelectionMode(QAbstractItemView::ExtendedSelection);
    //ui->listWidget->setDragEnabled(true);
    //ui->listWidget_2->viewport()->setAcceptDrops(true);
    //ui->listWidget_2->setDropIndicatorShown(true);
    ui->btnBack->setDisabled(true);
    ui->btnForward->setDisabled(true);
    ui->btnUp->setDisabled(true);
}

MainWindow::~MainWindow()
{
    delete ui;
}


void MainWindow::on_btnBack_clicked() {
    if (uploadBack.size()) {
        auto prevDir = uploadBack.top();
        uploadBack.pop();
        if (!uploadBack.size()) {
            ui->btnBack->setDisabled(true);
        }
        uploadForward.push(current_path);
        ui->btnForward->setEnabled(true);
        current_path = prevDir;
        ui->linePath->setText(current_path);

        QDir curr_dir(current_path);

        ui->listWidget->clear();
        for (const auto& item : curr_dir.entryInfoList()) {
            ui->listWidget->addItem(item.fileName());
        }

    }
}

void MainWindow::on_btnForward_clicked() {
    if (uploadForward.size()) {
        auto prevDir = uploadForward.top();
        uploadForward.pop();
        if (!uploadForward.size()) {
            ui->btnForward->setDisabled(true);
        }
        uploadBack.push(current_path);
        ui->btnBack->setEnabled(true);
        current_path = prevDir;
        ui->linePath->setText(current_path);

        QDir curr_dir(current_path);

        ui->listWidget->clear();
        for (const auto& item : curr_dir.entryInfoList()) {
            ui->listWidget->addItem(item.fileName());
        }

    }
}

void MainWindow::on_btnUp_clicked() {
    QString dst("..");
    QDir curr_dir(current_path);
    curr_dir.cd(dst);
    if (curr_dir.path() != current_path) {
        ui->btnBack->setEnabled(true);
        uploadBack.push(current_path);
    }
    current_path = curr_dir.path();
    ui->linePath->setText(current_path);
    ui->listWidget->clear();
    for (const auto& item : curr_dir.entryInfoList()) {
        ui->listWidget->addItem(item.fileName());
    }
}

void MainWindow::on_btnHome_clicked() {
    ui->btnUp->setEnabled(true);
    QString homepath(qgetenv("HOME"));
    if (current_path.length() && current_path != homepath) {
        ui->btnBack->setEnabled(true);
        uploadBack.push(current_path);
    }
    QDir dir(homepath);
    current_path = homepath;
    ui->linePath->setText(current_path);
    ui->listWidget->clear();
    for (const auto& item : dir.entryInfoList()) {
        ui->listWidget->addItem(item.fileName());
    }
}

void MainWindow::on_btnAdd_clicked() {
    for (const auto& item : ui->listWidget->selectedItems()) {
        QString fname = item->text();
        QString fullpath = current_path+"/"+fname;
        if (uploadset.find(fullpath) == uploadset.end()) {
            std::vector<fs::path> ld;
            ld.push_back(fs::path(fullpath.toStdString()));
            totalSize += EvaluateSize(ld, "");
            uploadset.insert(fullpath);
            ui->listWidget_2->addItem(fullpath);
        }
    }
    ui->lineTotal->setText(QString::number(totalSize));
}

void MainWindow::on_btnRemove_clicked() {
    for (QListWidgetItem *item : ui->listWidget_2->selectedItems()) {
        uploadset.erase(item->text());
        std::vector<fs::path> ld;
        ld.push_back(fs::path(item->text().toStdString()));
        totalSize -= EvaluateSize(ld, "");
        delete ui->listWidget_2->takeItem(ui->listWidget_2->row(item));
    }
    ui->lineTotal->setText(QString::number(totalSize));

}

void MainWindow::on_btmClear_clicked() {
    uploadset.clear();
    ui->listWidget_2->clear();
    totalSize = 0;
    ui->lineTotal->setText(QString::number(totalSize));
}

void MainWindow::on_listWidget_itemDoubleClicked(QListWidgetItem *item) {
//    QString dst = item->text();
//    QDir curr_dir(current_path);
//    curr_dir.cd(dst);
//    current_path = curr_dir.path();
//    ui->linePath->setText(current_path);
//    ui->listWidget->clear();
//    for (const auto& item : curr_dir.entryInfoList()) {
//        ui->listWidget->addItem(item.fileName());
//    }
}

void MainWindow::on_listWidget_itemActivated(QListWidgetItem *item) {
    QString dst = item->text();
    QDir curr_dir(current_path);
    curr_dir.cd(dst);
    if (curr_dir.path() != current_path) {
        uploadBack.push(current_path);
        ui->btnBack->setEnabled(true);
    }
    current_path = curr_dir.path();
    ui->linePath->setText(current_path);
    ui->listWidget->clear();
    for (const auto& item : curr_dir.entryInfoList()) {
        ui->listWidget->addItem(item.fileName());
    }
}

void MainWindow::on_btnPath_clicked()
{

}

void MainWindow::on_btnUpload_clicked() {
    struct visFuncs v;
    v.End1 = vis2;
    v.End2 = vis2;
    v.Next = vis3;
    v.Begin1 = vis2;
    v.Begin2 = vis2;
    v.SetField = vis1;
    QMessageBox* load = new QMessageBox();
    load->show();
    for (const auto& item : uploadset) {
        int res = UploadFile(item.toStdString(), "", 0, &v, 1600, 1);
    }
    load->close();
    delete load;
    uploadset.clear();
    ui->listWidget_2->clear();
    totalSize = 0;
    ui->lineTotal->setText(QString::number(totalSize));
}

void MainWindow::on_actionLogout_triggered() {
    QMessageBox::StandardButton btn =  QMessageBox::question(this, "Confirm action", "Are you sure to logout?", QMessageBox::Yes | QMessageBox::No);
    if (btn == QMessageBox::Yes) {
        ui->actionChange_User->setDisabled(true);
        ui->actionLogout->setDisabled(true);
        ui->actionUsername_some_thing->setText("username@someth.ing");
        ui->actionUsername_some_thing->setDisabled(true);
        std::fstream authFile(".auth");
        authFile.clear();
        authFile << false << "\n";
        authFile.close();
        ui->actionAuthorize->setEnabled(true);
        ui->actionUser_Options->setDisabled(true);
    }
}

void MainWindow::on_actionChange_User_triggered() {
    QMessageBox::StandardButton btn =  QMessageBox::question(this, "Confirm action", "Are you sure to change user?", QMessageBox::Yes | QMessageBox::No);
    if (btn == QMessageBox::Yes) {
        ui->actionChange_User->setDisabled(true);
        ui->actionLogout->setDisabled(true);
        ui->actionUser_Options->setDisabled(true);
        ui->actionAuthorize->setEnabled(true);
        ui->actionUsername_some_thing->setText("username@someth.ing");
        ui->actionUsername_some_thing->setDisabled(true);
        std::fstream authFile(".auth");
        authFile.clear();
        authFile << false << "\n";
        authFile.close();

        on_actionAuthorize_triggered();
    }

}

void MainWindow::on_actionAuthorize_triggered() {
    AuthWindow * auth = new AuthWindow();
    connect(auth, SIGNAL(AuthorizedLogin(QString)), this, SLOT(setAuthLogin(QString)));
    auth->setModal(true);
    auth->exec();
    disconnect(auth, SIGNAL(AuthorizedLogin(QString)), this, SLOT(setAuthLogin(QString)));
    delete auth;
}

void MainWindow::on_actionUser_Options_triggered()
{

}

void MainWindow::on_actionUsername_some_thing_triggered()
{

}

void MainWindow::setAuthLogin(const QString &auth_login) {
    ui->actionLogout->setEnabled(true);
    ui->actionChange_User->setEnabled(true);
    ui->actionUsername_some_thing->setEnabled(true);
    ui->actionUsername_some_thing->setText(auth_login);
    ui->actionUser_Options->setEnabled(true);
    ui->actionAuthorize->setDisabled(true);
}

void MainWindow::on_btnCd_clicked() {
    ui->btnHome->setEnabled(true);
    if (current_path.length() && current_path != ui->linePath->text()) {
        ui->btnBack->setEnabled(true);
        uploadBack.push(current_path);
    }
    current_path = ui->linePath->text();
    QDir dir(current_path);
    ui->linePath->setText(current_path);
    ui->listWidget->clear();
    for (const auto& item : dir.entryInfoList()) {
        ui->listWidget->addItem(item.fileName());
    }
}

unsigned long long MainWindow::EvaluateSize(std::vector<std::filesystem::__cxx11::path> &args, const std::string &start_path) {
    unsigned long long res = 0;
        for (auto& arg : args) {
            if (fs::exists(arg)) {
                if (fs::is_directory(arg)) {
                    for (const auto &entry_point : fs::recursive_directory_iterator(arg, fs::directory_options::skip_permission_denied)) {
                        if (fs::is_regular_file(entry_point.path())) {
                            res += entry_point.file_size();
                        }
                    }
                } else {
                    res += fs::file_size(arg);
                }
            } else {
                qDebug() << "No such file or directory :\n";
            }
        }
        return res;
}
