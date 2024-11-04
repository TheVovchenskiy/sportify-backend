# Sportify Backend

При изменениях вначале пересобираем:

```shell
export CONFIG_FILE=./config/config.example.yaml && \
export POSTGRES_ENV_FILE=config/postgres.example.env && \
make docker-compose-build
```

Сборка c помощью Makefile или вручную исполняя команды

```shell
export CONFIG_FILE=./config/config.example.yaml && \ 
export POSTGRES_ENV_FILE=config/postgres.example.env && \
make docker-compose-up
```

Для накатки миграций:

```shell
make migration-up
```

Если вдруг не сработало, попробуйте:

```shell
make migration-up-reserve
```

Заполнения бд (нужно только один раз) можно вручную исполнить из IDE [sql запрос](sportify/db/fill.sql) или
вот так:

```shell
make fill-db
```

Если вы видите что-то такое "duplicate key value violates unique constraint", то скорее всего бд уже заполнена.

Все, вы прекрасны)

Посмотреть логи только backend контейнера:
```shell
make docker-compose-logs
```

### Для разработчиков

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

Для запуска чисто бэка(без бд) на локалке:

```shell
cd sportify && go run . --config-path=<path_to_config_dir>
```

Тут `path_to_config_dir` - путь до директории с конфигом, можно указывать несколько. При этом сам конфиг обязательно должен называться `config.yaml`.
