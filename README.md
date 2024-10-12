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

Для запуска бэка на порту 8080: 
```shell
cd sportify && go run .
```