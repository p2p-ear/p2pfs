# Техническое Задание
## Введение
### Цель
Удобное облачное хранение данных с высокой степенью защищенности и дублирования информации.
### Краткая сводка возможностей
- Распределенное хранение файлов
- Защита от поломок, поддерживание коэффициента репликации
- Кусочное хранение больших файлов
- Удобное управление файлами с сайта
## Детальные требования
### Варианты применения
- Персональное облако (без шифрования) - безопасное хранение больших данных
- Многопользовательский режим: авторизация, шифрование
- Аренда (и сдача в аренду) дискового пространства - монетизация
### Функциональные требования
- Скорость
- Масштабируемость
- Надежность
- Максимальная децентрализация
### Нефункциональные требования
- Удобное взаимодействие с программой из коммандной строки
- Удобный веб-интерфейс
### Дополнительные требования
- Визуализация статистики
- Локальный GUI-интерфейс
## Средства реализации
### Golang:
- сетевая структура
- сетевые алгоритмы: репликация, gossip
### Python:
- Взаимодействие с базой данных пользователей
- Бэкенд для сайта (Flask)
### C++:
- Шифрование
- Сжатие
- Шардинг
### HTML & CSS:
- Фронтэнд
 