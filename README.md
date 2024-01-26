# Приложение EWallet

### Запуск
```shell
$ git clone https://github.com/timohahaa/ewallet.git
$ cd ewallet
$ docker-compose up
```
__Перед первым запуском приложения необходимо накинуть миграцию `up.sql` из папки `migrations` в контейнере с базой данных.__
После все данные будут персистентными - создастся отдельный docker-volume и данные не будут теряться при перезапуске приложения.
Так же создастся отдельный docker volume под логи.

__*Не забудьте указать параметры подключения к базе и настройки приложения в файлах `config.yaml` (см. папку `config`) и `.env` (см. файл `.env.example`)*__

### Стэк и использованные библиотеки
- Go 1.21
- PostgreSQL (14 версия)
- docker + docker-compose
- https://github.com/ilyakaznacheev/cleanenv - библиотека для работы с конфигами
- https://github.com/labstack/echo - http-роутер
- https://github.com/sirupsen/logrus - логирование
- https://github.com/timohahaa/postgres - __*Самописная*__ библиотека для работы с PostgreSQL (pgx - драйвер И squirrel - sql-builder)
- https://github.com/google/uuid - пакет для работы с UUID

### Выполнение требований
##### Безопасность: в приложении не должно быть уязвимостей, позволяющих произвольно менять данные в базе.
 - Достигается за счет структуры запросов к базе/апи приложения + валидирования sql-иньекций на уровне билдера запросов и интерфейса драйвера базы
##### Персистентность: данные и изменения не должны «теряться» при перезапуске приложения.
 - Достаточно накатить миграцию только один раз - при первом запуске приложения - далее создастся docker-volume для контейнера с базой данных и данные не будут теряться при перезапуске/падении приложения

 ### Как протестировать API?
 Лично я рекомендую Postman
 Но вот список curl-ов для случая, если нет возможности использовать Postman:
 (здесь сервер запущен на порту 8080)

 Эндпоинт - POST /api/v1/wallet
 ```shell
 $ curl --location --request POST 'http://localhost:8080/api/v1/wallet' \
--header 'Content-Type: application/json'
 ```

 Эндпоинт - POST /api/v1/wallet/{walletId}/send
 (создайте перед этим два кошелька и замените указанные в запросе на свои)
```shell
$ curl --location 'http://localhost:8080/api/v1/wallet/05bb88df-eef6-4b6e-b024-a3d9d7448e6c/send' \
--header 'Content-Type: application/json' \
--data '{
    "to": "05bb88df-eef6-4b6e-b024-a3d9d7448e6c",
    "amount": 25.0
}'
```

Эндпоинт – GET /api/v1/wallet/{walletId}/history
```shell
$ curl --location 'http://localhost:8080/api/v1/wallet/05bb88df-eef6-4b6e-b024-a3d9d7448e6c/history' \
--header 'Content-Type: application/json'
```

Эндпоинт – GET /api/v1/wallet/{walletId}
```shell
$ curl --location 'http://localhost:8080/api/v1/wallet/cdb494a1-7819-4dec-9ed6-0f7a88884da9' \
--header 'Content-Type: application/json'
```