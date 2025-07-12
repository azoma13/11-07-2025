# archiving-service

## Описание проекта:
Сервис, который скачивает файлы с открытого доступа в интернете и запаковывает их в zip архив. 

## Стек:
    Язык сервиса: Go.

## Сервис предоставляет следующие конечные точки API:
- создание задачи на создание архива;
- добавление ссылки на файл, который будет в архиве для каждой задачи;
- получение статуса задачи. При заполнении задачи максимальным количеством файлов, указанный в `.env` параметра `APP_MAX_NUM_FILES`, данная точка возвращает совместно со статусом ссылку на архив.

## Заполнение `.env` файл
Пример `.env` файла:
```
APP_NAME=archiving-service                  // имя проекта
APP_VERSION=1.0.0                           // версия проекта
APP_MAX_NUM_TASKS=3                         // максимальное количество задач
APP_MAX_NUM_FILES=3                         // максимальное количество файлов в задаче
APP_ALLOWED_FILE_EXTENSIONS=.jpeg,.pdf      // формат файлов поддерживающий на скачивание

HTTP_PORT=8080                              // порт сервера
```

## Инструкция по запуску проекта
1. Описать файл `.env` согласно примеру `.env.example`
2. Запустить сервер командой: `make run`

## Примеры
Примеры возможных запросов:
- [создание задачи](#create)
- [добавление ссылки на файл](#add-file)
- [получение статуса задачи](#status)

### Создание задачи <a name="create"></a>
Запрос:
```h
POST /api/v1/task/create HTTP/1.1
Host: localhost:8080
```
Ответ: `201 Created`
```json
{
    "task_id": 1
}
```

### Добавление ссылки на файл <a name="add-file"></a>
Запрос:
```h
POST /api/v1/task/add-file HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 132

{
    "task_id": 1,
    "url_file": "https://mrroot.pro/wp-content/uploads/2023/10/black-hat-go.pdf?ysclid=md0h6kdfp5929139483"
}
```
Ответ: `204 No Content`

### Получение статуса задачи <a name="status"></a>
Запрос:
```h
GET /api/v1/task/status?task_id=1 HTTP/1.1
Host: localhost:8080
```
Ответ при пустом архиве: `200 OK`
```json
{
    "status_content": "Empty"
}
```
Ответ при не пустом и не полном архиве: `200 OK`
```json
{
    "status_content": "In progress"
}
```
Ответ при полном архиве: `200 OK`
```json
{
    "status_content": "Complete",
    "archiving_url": "/home/example/example/archiving-service/archives/task_1.zip"
}
```