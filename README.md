# BannerCounter
![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)
![fasthttp](https://img.shields.io/badge/fasthttp-fast-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?logo=postgresql&logoColor=white)

BannerCounter — это сервис для подсчёта и получения статистики показов баннеров по минутам.

# Переменные окружения

Сервис настраивается через следующие переменные окружения:
| Переменная   | Описание |
|--------------|----------|
| `INTERVAL`   | Интервал сохранения накопленных данных в базу. Формат Go Duration (например: `5s`, `1m`). |
| `LISTEN_PORT`| Порт, на котором сервис будет слушать HTTP-запросы (по умолчанию `5454`). |
| `HOST`       | Хост PostgreSQL. |
| `PORT`       | Порт PostgreSQL (по умолчанию `5432`). |
| `USER`       | Пользователь PostgreSQL. |
| `PASSWORD`   | Пароль PostgreSQL. |
| `DB_NAME`    | Имя базы данных. |


Пример .env файла:
``` bash
INTERVAL=5s
LISTEN_PORT=5454
HOST=postgres
PORT=5432
USER=postgres
PASSWORD=postgres
DB_NAME=postgres
```

# Запуск через Docker Compose

Для удобства есть docker-compose.yml, который запускает всё сразу: сервис и базу данных PostgreSQL.

## Подготовка
``` bash
task build
```
или

``` bash
docker build -t bannercounter:latest -f ./Dockerfile .
```
## Запуск
``` bash
docker compose up
```

## После запуска

API сервиса будет доступно на http://localhost:5454

PostgreSQL будет доступен на порту 5432 внутри Docker сети (или на хосте, если проброшен).

# Основные эндпоинты

- Добавление показов (GET)
    `GET /counter/{bannerID}` — увеличивает счётчик показов баннера.
    Пример запроса:
    ```bash
    curl -X GET "http://localhost:5454/counter/1"
    ```

- Получение статистики (POST)
    `POST /stats/{bannerID}` — выдаёт статистику за заданный промежуток времени.
    Пример запроса:
    ``` bash
    curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"from": "2025-08-15T14:00:00", "to": "2025-08-15T19:00:00"}' \
    "http://localhost:5454/stats/1"
    ```

    Пример ответа:
    ``` json
    {
    "stats": 
        [
            {"ts": "2025-08-15T14:00:00", "v": 4},
            {"ts": "2025-08-15T14:01:00", "v": 2}
        ]
    }
    ```