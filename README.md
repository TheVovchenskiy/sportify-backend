# Move-Life Backend (ex. sportify)

Это backend проекта Move-Life - инструмента для упрощения поиска и 
организации спортивных игр. Он состоит из двух частей сайта и интеграции с телеграмм.

### [Ссылка на сайт](https://move-life.ru/events)
### [Ссылка на бота](https://t.me/movelife_ond_bot)
### [Репозиторий фронтенд](https://github.com/DriverOnLips/sportify-frontend)

## Бот в телеграмм
Взаимодействие с телеграмм происходит через бота на python директория /bot. 
Его зона ответственности все что связано с телеграммом: доставка сообщений, 
изменение сообщений, подтверждение авторизации через телеграмм, уведомления и т.д.

## Sportify сервис
Основной сервис на Golang директория /sportify, точка доступа к бд Postgres
и к хранилищу файлов, апи для фронтенда, периодические задачи,
походы в сторонние апи.

## Nginx proxy на входе
Служит для цели ssl терминации и как легковесная точка входа в наши сервисы.

# Команды и инструкции для разработки

### Уважаемые фронты, выполните эту команду один раз (дальше не надо будет). Чтобы на локалке сгенерить сертификаты.
```shell
make gen-cert
```

### При изменениях вначале пересобираем:

```shell
export CONFIG_FILE=./config/config.example.yaml && \
export POSTGRES_ENV_FILE=config/postgres.example.env && \
make docker-compose-build
```

### Сборка c помощью Makefile или вручную исполняя команды

```shell
export CONFIG_FILE=./config/config.example.yaml && \ 
export POSTGRES_ENV_FILE=config/postgres.example.env && \
make docker-compose-up
```

### Для накатки миграций:

```shell
make migration-up
```

### Если вдруг не сработало, попробуйте:

```shell
make migration-up-reserve
```

### Для заполнения бд (нужно только один раз) можно вручную исполнить из IDE [sql запрос](sportify/db/fill.sql) или
вот так:

```shell
make fill-db
```

Если вы видите что-то такое "duplicate key value violates unique constraint", то скорее всего бд уже заполнена.

Все, вы прекрасны)

### Посмотреть логи только backend контейнера:
```shell
make docker-compose-logs
```
