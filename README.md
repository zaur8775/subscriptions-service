# Stack

* Go 1.26
* PostgreSQL 15
* Docker & Docker Compose
* Swagger (swaggo)
* golang-migrate
* pgx
* net/http

---

# Запуск проекта

## Запуск через Docker

```bash
make docker-up
```

или:

```bash
docker compose up --build
```

---

# Swagger документация

После запуска доступна по адресу:

```text
http://localhost:8080/swagger/index.html
```

---

# Миграции

## Применить миграции

```bash
make migrate-up
```

## Откатить миграции

```bash
make migrate-down
```

## Создать новую миграцию

```bash
make migrate-create name=create_table
```

---

# Makefile команды
            
`make run`                       запуск приложения локально
`make build`                     сборка приложения
`make swagger`                   генерация swagger документации
`make docker-up`                 запуск docker compose
`make docker-down`               остановка docker compose
`make migrate-up`                применение миграций       
`make migrate-down`              откат миграций            
`make migrate-create name=...`   создание миграции
