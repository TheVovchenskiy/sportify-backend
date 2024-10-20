# Sportify Backend

Установка го:
```shell
wget --directory-prefix=bin https://go.dev/dl/go1.23.2.linux-amd64.tar.gz && \
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf bin/go1.23.2.linux-amd64.tar.gz && \
export PATH=$PATH:/usr/local/go/bin
```

Для проверки версии:
```shell
go version
```
Вывод примерно такой: go version go1.23.0 linux/amd64

Установка всего тулчейна:
```shell
make toolchain
```

Сборка c помощью Makefile или вручную исполняя команды
```shell
export CONFIG_FILE=config/config.example.yaml && \
export POSTGRES_ENV_FILE=../config/postgres.example.env && \
make docker-compose-up
```

Для накатки миграций:
```shell
migrate -database postgres://postgres:postgres@localhost:5432/sportify?sslmode=disable -path ./sportify/db/migrations up
```

Для заполнения бд можно вручную исполнить из IDE [sql запрос](sportify/db/fill.sql) или 
зайти в контейнер бд (docker exec -it deploy-postgres-1 psql -U postgres -d sportify),
а там уже запустить скрипт.

Для запуска чисто бэка(без бд) на локалке:
```shell
export CONFIG_FILE=config/config.example.yaml && \
cd sportify && go run . -configfile=config.example
```

Для пересборки проекта после изменений:
```shell
make docker-compose-build
```