# icpo-backend

REST API сервис на Go для ведения картотеки зарегистрированных преступников Интерпола. Обеспечивает регистрацию, поиск по произвольному подмножеству признаков, архивацию «завязавших» и удаление записей только после подтверждённой смерти.

## Быстрый старт

Требования: Go и Docker.

```bash
# 1. Создать файл с переменными окружения
cp .env.example .env
 
# 2. Поднять PostgreSQL
docker compose up -d postgres
 
# 3. Запустить сервер (миграции применятся автоматически)
go run ./cmd/api
```

Проверка, что всё работает:

```bash
curl http://localhost:8080/api/v1/health
# {"status":"ok"}
```

Документация API (Swagger UI): <http://localhost:8080/swagger/>

Остановить сервер — `Ctrl+C` (штатное завершение, до 10 секунд). Остановить базу — `docker compose down` (данные сохранятся в volume `postgres_data`; чтобы удалить и их — `docker compose down -v`).

Сборка бинарника:

```bash
go build -o bin/api ./cmd/api
```


## Конфигурация

Значения берутся из переменных окружения. Если в рабочей директории есть файл `.env`, он подгружается автоматически. Образец — `.env.example`.

| Переменная          | По умолчанию | Описание                                  |
|---------------------|--------------|-------------------------------------------|
| `HTTP_ADDR`         | `:8080`      | Адрес, на котором слушает HTTP-сервер     |
| `POSTGRES_HOST`     | `localhost`  | Хост PostgreSQL                           |
| `POSTGRES_PORT`     | `5432`       | Порт PostgreSQL                           |
| `POSTGRES_DB`       | `postgres`   | Имя базы данных                           |
| `POSTGRES_USER`     | `postgres`   | Пользователь                              |
| `POSTGRES_PASSWORD` | _(пусто)_    | Пароль                                    |
| `POSTGRES_SSLMODE`  | `disable`    | Режим SSL (`disable`, `require`, ...)     |
| `JWT_SECRET`        | _(обязательна)_ | Секрет для подписи access-токенов (HS256); без него сервер не запустится |
| `JWT_ACCESS_TTL`    | `15m`        | Время жизни access-токена                 |
| `JWT_REFRESH_TTL`   | `720h`       | Время жизни refresh-токена                |

Строка подключения собирается из `POSTGRES_*`:
`postgres://USER:PASSWORD@HOST:PORT/DB?sslmode=SSLMODE`.

Эти же `POSTGRES_*` использует `docker-compose.yaml` для создания пользователя и базы в контейнере.

Переменные `GOOSE_*` из `.env.example` нужны только для ручного запуска CLI `goose` (см. раздел «Миграции»). Сам сервер их не читает.

## Миграции

Миграции лежат в каталоге `migrations/` в виде SQL-файлов и подключаются в бинарник директивой `//go:embed *.sql`. При каждом старте сервер выполняет `goose up`: применяет все ещё не применённые миграции, а в служебной таблице `goose_db_version` хранится, какие уже выполнены.

Важно:

- Сервер умеет только **применять** миграции (`up`). Откат (`down`) выполняется вручную через CLI `goose`.
- Файлы встроены в бинарник при сборке, поэтому после добавления файла нужно пересобрать или перезапустить через `go run` — старый бинарник новую миграцию не увидит.
### Установка goose CLI

Ставьте ту же версию, что в `go.mod`:

```bash
go install github.com/pressly/goose/v3/cmd/goose@v3.28.0
```

Убедитесь, что `$(go env GOPATH)/bin` есть в `PATH`.

### Подключение CLI к базе

CLI читает параметры из переменных окружения `GOOSE_DRIVER`, `GOOSE_DBSTRING` и `GOOSE_MIGRATION_DIR`. Проще всего загрузить их из `.env`:

```bash
set -a; source .env; set +a
```

После этого команды можно запускать без флагов. Либо укажите всё явно:

```bash
goose -dir migrations postgres "postgres://admin:admin@localhost:5432/icpo_db?sslmode=disable" status
```

### Добавление миграции

1. Создайте файл (имя получит префикс с временем):
```bash
   goose -dir migrations create create_users_table sql
```

Появится файл вида `migrations/20261009120000_create_users_table.sql` с заготовкой.

2. Заполните разделы `Up` (что сделать) и `Down` (как это откатить):
```sql
   -- +goose Up
   CREATE TABLE users (
       id         BIGSERIAL PRIMARY KEY,
       email      TEXT        NOT NULL UNIQUE,
       created_at TIMESTAMPTZ NOT NULL DEFAULT now()
   );
 
   -- +goose Down
   DROP TABLE users;
```

Если в миграции есть блоки с `;` внутри (функции, триггеры, `DO $$ ... $$`), оборачивайте их в `-- +goose StatementBegin` / `-- +goose StatementEnd`:

```sql
   -- +goose Up
   -- +goose StatementBegin
   CREATE FUNCTION set_updated_at() RETURNS trigger AS $$
   BEGIN
       NEW.updated_at = now();
       RETURN NEW;
   END;
   $$ LANGUAGE plpgsql;
   -- +goose StatementEnd
 
   -- +goose Down
   DROP FUNCTION set_updated_at();
```

3. Примените к базе одним из способов:
    - **Запустить сервер** — `go run ./cmd/api`; миграция применится на старте.
    - **Через CLI**, не запуская сервер: `goose -dir migrations up`.
4. Проверьте результат:
```bash
   goose -dir migrations status
```

### Удаление (откат) миграции

Всё зависит от того, применена ли миграция.

**Миграция ещё нигде не применялась** (только создали файл) — просто удалите файл из `migrations/`.

**Миграция уже применена к базе** — сначала откатите её, потом удаляйте файл:

```bash
# откатить последнюю применённую миграцию
goose -dir migrations down
 
# проверить, что она теперь в статусе Pending
goose -dir migrations status
 
# после этого удалить файл
rm migrations/20261009120000_create_users_table.sql
```

Если сначала удалить файл, а потом откатывать — `goose down` не найдёт SQL для отката, и изменения в базе останутся.

> Откат с удалением файла подходит для локальной разработки. Если миграция уже применена на общем/боевом окружении, не удаляйте её: напишите **новую** миграцию, которая отменяет изменения (например, `DROP TABLE users;` в разделе `Up`), и примените её как обычно.

### Шпаргалка по командам

Все команды выполняются из корня проекта с загруженным `.env` (см. выше).

| Команда                                    | Что делает                                                    |
|--------------------------------------------|---------------------------------------------------------------|
| `goose -dir migrations status`             | Показать, какие миграции применены, а какие нет               |
| `goose -dir migrations up`                 | Применить все неприменённые миграции                          |
| `goose -dir migrations up-by-one`          | Применить только следующую миграцию                           |
| `goose -dir migrations up-to 20261009120000` | Применить миграции до указанной версии включительно         |
| `goose -dir migrations down`               | Откатить последнюю применённую миграцию                       |
| `goose -dir migrations down-to 20261008163116` | Откатить до указанной версии (она остаётся применённой)   |
| `goose -dir migrations redo`               | Откатить последнюю миграцию и применить её заново             |
| `goose -dir migrations reset`              | Откатить **все** миграции (осторожно: данные будут потеряны)  |
| `goose -dir migrations version`            | Показать текущую версию схемы                                 |
| `goose -dir migrations create NAME sql`    | Создать новый SQL-файл миграции                               |

### Полный сброс локальной базы

Если нужно начать с чистого листа, самый быстрый способ — удалить volume с данными и поднять базу заново; сервер применит все миграции с нуля при старте:

```bash
docker compose down -v
docker compose up -d postgres
go run ./cmd/api
```