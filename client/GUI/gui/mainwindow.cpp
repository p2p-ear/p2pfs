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
    manager = new QNetworkAccessManager;
    QObject::connect(
                manager,
                &QNetworkAccessManager::finished,
                // Лямбда функция - обработчик HTTP ответа
                [=](QNetworkReply *reply) {
        QString responce;
        ui->textBrowser->clear();
        // Обработка ошибок
        if (reply->error()) {
            responce += QString("Error %1").arg(reply->errorString())+"\n";
            reply->deleteLater();
        }

        // Вывод заголовков
//        for (auto &i:reply->rawHeaderPairs()) {
//            QString str;
//            responce += str.sprintf(
//                            "%40s: %s",
//                            i.first.data(),
//                            i.second.data());
//        }

        // Вывод стандартного заголовка
        responce += reply->header(QNetworkRequest::ContentTypeHeader).toString()+"\n";

        // Тело ответа в формате JSON
        QByteArray responseData = reply->readAll();
        QJsonDocument doc(QJsonDocument::fromJson(responseData));
        responce += doc.toJson();
        QJsonObject rep = doc.object();
        if (rep["status"].toInt() != 0) {
            QMessageBox::information(this, "Error", rep["message"].toString());
        } else if (rep["email"] != Login) {
            QMessageBox::warning(this, "Authoriztion", "Authorization failed");
        } else {}

        ui->textBrowser->setText(responce);
        // Delete garbage && Exit
        reply->deleteLater();
    });


    ui->setupUi(this);
    ui->actionUsername_some_thing->setDisabled(true);
    ui->actionLogout->setDisabled(true);
    ui->actionChange_User->setDisabled(true);
    ui->actionUser_Options->setDisabled(true);
    ui->tableWidget->setSelectionMode(QAbstractItemView::ExtendedSelection);
    ui->tableWidget->setSelectionBehavior(QAbstractItemView::SelectRows);
    //ui->listWidget->setDragEnabled(true);
    //ui->listWidget_2->viewport()->setAcceptDrops(true);
    //ui->listWidget_2->setDropIndicatorShown(true);
    ui->btnBack->setDisabled(true);
    ui->btnForward->setDisabled(true);
    ui->btnUp->setDisabled(true);
    ui->tableWidget->resizeColumnsToContents();
    ui->tableWidget_2->resizeColumnsToContents();
    FS.Load();

    //non-func buttons
    ui->btnBack2->setDisabled(true);
    ui->btnForward2->setDisabled(true);
    ui->btnHome2->setDisabled(true);
    ui->btnUp2->setDisabled(true);
    ui->btnPath2->setDisabled(true);
    ui->btnPath->setDisabled(true);
    ui->btnUpload2->setDisabled(true);
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

        updateTable(current_path);
        on_btmClear_clicked();
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

        updateTable(current_path);
        on_btmClear_clicked();
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

    updateTable(current_path);
    on_btmClear_clicked();
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
    updateTable(current_path);
    on_btmClear_clicked();
}

void MainWindow::on_btnAdd_clicked() {
    for (const auto& item : ui->tableWidget->selectedItems()) {
        if (item->column() == 0) {
            QString fname = item->text();
            QString fullpath;
            if (current_path[current_path.length()-1] == '/') {
                fullpath = current_path+fname;
            } else {
                fullpath = current_path+"/"+fname;
            }
            if (uploadset.find(fullpath) == uploadset.end()) {
                updateTable2(fullpath);
            }
        }
    }
    ui->lineTotal->setText(QString::number(totalSize));
}

void MainWindow::on_btnRemove_clicked() {
    for (QTableWidgetItem *item : ui->tableWidget_2->selectedItems()) {
        if (item != nullptr && ui->tableWidget_2->column(item) == 0) {
            uploadset.erase(item->text());
            std::vector<fs::path> ld;
            ld.push_back(fs::path(item->text().toStdString()));
            totalSize -= EvaluateSize(ld, "");
            ui->tableWidget_2->removeRow(ui->tableWidget_2->row(item));
        }
    }
    ui->lineTotal->setText(QString::number(totalSize));

}

void MainWindow::on_btmClear_clicked() {
    uploadset.clear();
    ui->tableWidget_2->setRowCount(0);
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
    std::string ip = "198.172.0.1:9000";
    unsigned long ringsz = 1000;
//    for (const auto& item : uploadset) {
//        int res = UploadFile(item.toStdString(), "", 0, &v, 1600, 1, ip, ringsz);
//    }
    QJsonArray arr;
    for (int i = 0; i < ui->tableWidget_2->rowCount(); i++) {
        UploadFile(ui->tableWidget_2->item(i, 0)->text().toStdString(), "", 0, &v, 1600, 1, ip, ringsz);
        QJsonObject obj;
        obj.insert(T_NAME, ui->tableWidget_2->item(i, 0)->text());
        obj.insert(T_ISDIR, ui->tableWidget_2->item(i, 1)->text() == "dir" ? true : false);
        obj.insert(T_SIZE, ui->tableWidget_2->item(i, 2)->text());
        if (ui->tableWidget_2->item(i, 1)->text() == "dir") {
            obj.insert(T_CHILD, QJsonArray());
        }
        arr.append(obj);
    }
    //arr -- array of jsons, pathToLoad -- path on MyDisk
    QString pathToLoad = ui->linePath->text();

    load->close();
    delete load;
    uploadset.clear();
    ui->tableWidget_2->clear();
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
        is_authorised = false;
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
        is_authorised = false;
        on_actionAuthorize_triggered();
    }

}

void MainWindow::on_actionAuthorize_triggered() {
    AuthWindow * auth = new AuthWindow();
    connect(auth, SIGNAL(AuthorizedLogin(QString, QString, QString)), this, SLOT(setAuthLogin(QString, QString, QString)));
    auth->setModal(true);
    auth->exec();
    disconnect(auth, SIGNAL(AuthorizedLogin(QString, QString, QString)), this, SLOT(setAuthLogin(QString, QString, QString)));
    delete auth;
}

void MainWindow::on_actionUser_Options_triggered()
{

}

void MainWindow::on_actionUsername_some_thing_triggered()
{

}

void MainWindow::setAuthLogin(const QString &auth_login, const QString & auth_pass, const QString& auth_JWT) {
    Login = auth_login;
    Password = auth_pass;
    JWT = auth_JWT;
    ui->actionLogout->setEnabled(true);
    ui->actionChange_User->setEnabled(true);
    ui->actionUsername_some_thing->setEnabled(true);
    ui->actionUsername_some_thing->setText(auth_login);
    ui->actionUser_Options->setEnabled(true);
    ui->actionAuthorize->setDisabled(true);
    is_authorised = true;
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
    updateTable(current_path);
    on_btmClear_clicked();
}

void MainWindow::updateTable(const QString& currentPath) {
    ui->tableWidget->setRowCount(0);
    QDir dir(currentPath);
    for (const auto& item : dir.entryInfoList()) {
        if (item.fileName() == ".") {
            continue;
        }
        ui->tableWidget->insertRow(ui->tableWidget->rowCount());
        QTableWidgetItem * Name = new QTableWidgetItem(item.fileName());
        QTableWidgetItem * Type = new QTableWidgetItem(item.isDir() ? "dir" : "file");
        QTableWidgetItem * Size = new QTableWidgetItem(QString::number(item.size()));
        ui->tableWidget->setItem(ui->tableWidget->rowCount()-1, 0, Name);
        ui->tableWidget->setItem(ui->tableWidget->rowCount()-1, 1, Type);
        ui->tableWidget->setItem(ui->tableWidget->rowCount()-1, 2, Size);
        ui->tableWidget->resizeColumnsToContents();
    }
}

void MainWindow::updateTable2(const QString & fullpath) {
    std::vector<fs::path> ld;
    ld.push_back(fs::path(fullpath.toStdString()));
    unsigned long long size = EvaluateSize(ld, "");
    totalSize += size;
    uploadset.insert(fullpath);
    ui->tableWidget_2->insertRow(ui->tableWidget_2->rowCount());
    QTableWidgetItem * Name = new QTableWidgetItem(fullpath);
    QTableWidgetItem * Type = new QTableWidgetItem(QFileInfo(fullpath).isDir() ? "dir" : "file");
    QTableWidgetItem * Size = new QTableWidgetItem(QString::number(size));

    ui->tableWidget_2->setItem(ui->tableWidget_2->rowCount()-1, 0, Name);
    ui->tableWidget_2->setItem(ui->tableWidget_2->rowCount()-1, 1, Type);
    ui->tableWidget_2->setItem(ui->tableWidget_2->rowCount()-1, 2, Size);
    ui->tableWidget_2->resizeColumnsToContents();
}

void MainWindow::updateTable3(const std::vector<MDfile> &src) {
    ui->tableWidget_3->setRowCount(0);
    for (const auto& item : src) {
        ui->tableWidget_3->insertRow(ui->tableWidget_3->rowCount());
        QTableWidgetItem * Name = new QTableWidgetItem(item.Name);
        QTableWidgetItem * Type = new QTableWidgetItem(item.isDir ? "dir" : "file");
        QTableWidgetItem * Size = new QTableWidgetItem(QString::number(item.Size));
        ui->tableWidget_3->setItem(ui->tableWidget_3->rowCount()-1, 0, Name);
        ui->tableWidget_3->setItem(ui->tableWidget_3->rowCount()-1, 1, Type);
        ui->tableWidget_3->setItem(ui->tableWidget_3->rowCount()-1, 2, Size);
        ui->tableWidget_3->resizeColumnsToContents();
    }
}

void MainWindow::updateTable4(const MDfile & file) {
    downloadset.insert(file.Name);
    ui->tableWidget_4->insertRow(ui->tableWidget_4->rowCount());
    QTableWidgetItem * Name = new QTableWidgetItem(file.Name);
    QTableWidgetItem * Type = new QTableWidgetItem(file.isDir ? "dir" : "file");
    QTableWidgetItem * Size = new QTableWidgetItem(QString::number(file.Size));

    ui->tableWidget_4->setItem(ui->tableWidget_4->rowCount()-1, 0, Name);
    ui->tableWidget_4->setItem(ui->tableWidget_4->rowCount()-1, 1, Type);
    ui->tableWidget_4->setItem(ui->tableWidget_4->rowCount()-1, 2, Size);
    ui->tableWidget_4->resizeColumnsToContents();
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

void MainWindow::on_tableWidget_itemActivated(QTableWidgetItem *item) {
    QString dst = ui->tableWidget->item(item->row(), 0)->text();
    QDir curr_dir(current_path);
    curr_dir.cd(dst);
    if (curr_dir.path() != current_path) {
        uploadBack.push(current_path);
        ui->btnBack->setEnabled(true);
    }
    current_path = curr_dir.path();
    ui->linePath->setText(current_path);
    updateTable(current_path);
}

void MainWindow::on_btnCd_2_clicked() {
    if (FS.Cd(ui->linePath_2->text())) {
        updateTable3(FS.Ls());
        on_btmClear2_clicked();
    }
}

void MainWindow::on_btnAdd2_clicked() {
    for (const auto& item : ui->tableWidget_3->selectedItems()) {
        if (item->column() == 0) {
            MDfile file;
            QString fname = item->text();
            QString fullpath;
            QString cpath = FS.GetCurrPath();
            if (cpath[cpath.length()-1] == '/') {
                fullpath = cpath+fname;
            } else {
                fullpath = cpath+"/"+fname;
            }
            file.Name = fullpath;
            file.Size = ui->tableWidget_3->item(item->row(), 2)->text().toUInt();
            file.isDir = ui->tableWidget_3->item(item->row(), 2)->text() == "dir";
            if (downloadset.find(fullpath) == downloadset.end()) {
                updateTable4(file);
            }
        }
    }
    //ui->lineTotal->setText(QString::number(totalSize));
}

void MainWindow::on_btnRemove2_clicked() {
    for (QTableWidgetItem *item : ui->tableWidget_4->selectedItems()) {
        if (item != nullptr && ui->tableWidget_4->row(item) == 0) {
            downloadset.erase(item->text());
            ui->tableWidget_4->removeRow(ui->tableWidget_4->row(item));
        }
    }
}



void MainWindow::on_btmClear2_clicked() {
    downloadset.clear();
    ui->tableWidget_4->setRowCount(0);
}

void MainWindow::on_pushButton_clicked() {

}

void MainWindow::on_pushButton_2_clicked() {

}

void MainWindow::on_btnAddCoins_clicked() {
    if (is_authorised) {
        QJsonObject jBody;
        jBody.insert("value", ui->lineEdit->text().toInt());
        MakeReqRequest(jBody, 2);
    }
}

void MainWindow::on_btnUodateJson_clicked() {
    if (is_authorised) {
        QJsonObject jBody;
        //jBody.insert("Null", "Null");
        MakeReqRequest(jBody, 3);
    }
}

int MainWindow::MakeReqRequest(QJsonObject &body, int type) {
    QJsonObject jObj;
    jObj.insert("email", Login);
    jObj.insert("pass", Password);
    jObj.insert("JWT", "a");//JWT);
    jObj.insert("type", type);
    jObj.insert("body", body);
    qDebug() << jObj.keys().size();
    QJsonDocument jDoc(jObj);
    QNetworkRequest req;
    req.setUrl(QUrl(addressRequest));
    qDebug() << jDoc.toJson();
    req.setRawHeader("Content-Type","application/json");
    manager->post(req, jDoc.toJson());
    return 1;
}

void MainWindow::on_btnGetCoins_clicked() {
    if (is_authorised) {
        QJsonObject jBody;
        //jBody.insert("Null", "Null");
        MakeReqRequest(jBody, 4);
    }
}
