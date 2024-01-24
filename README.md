# Инструмент "SQL-мигратор"

[![Test](https://github.com/XanderKon/sql-migrator-otus/actions/workflows/actions.yml/badge.svg)](https://github.com/XanderKon/sql-migrator-otus/actions/workflows/actions.yml)
[![coverage](https://raw.githubusercontent.com/XanderKon/sql-migrator-otus/badges/.badges/main/coverage.svg)](/.github/.testcoverage.yml)

## Технологии

go ~1.21

Поддерживаемый драйвер БД: `PSQL ^14`

## Общее описание

Аналог инструментов, приведенных в секции "Database schema migration"
[awesome-go](https://github.com/avelino/awesome-go).

Утилита, работающая с миграциями, написанными на Go или представленными в виде SQL-файлов.

Позволяет:

- генерировать шаблон миграции;
- применять миграции;
- откатывать миграции.

## Конфигурация

Задаётся в файле `configs/config.yml`

| Параметр   | Описание                       | Возможные значения                                                       |
| ---------- | ------------------------------ | ------------------------------------------------------------------------ |
| dsn        | Строка подключения к БД        | postgresql://app:!ChangeMe!@pgsql:5432/app?serverVersion=15&charset=utf8 |
| dir        | Директория для хранения файлов | migrations                                                               |
| type       | Тип генерируемых файлов        | go \| sql                                                                |
| table_name | Название таблицы в БД          | migrations                                                               |

В конфигурации можно использовать переменные окружения, тогда в качестве значения нужно использовать специальную нотацию: `${ENV_VAR}` или `$ENV_VAR`.

Пример:

`configs/config.yml`

```yml
dsn: "$DB_DSN"
dir: ./migrations
type: sql
table_name: migrations
```

`DB_DSN=SOME_ENV_HERE ./bin/gomigrator`

## Использование

### Как CLI-утилита

Устанавливаем

```bash
go install github.com/XanderKon/sql-migrator-otus/cmd/gomigrator@latest
```

Создаем файл конфигурации в нужном месте, следующего содержания:

```yml
migrator:
  dsn: postgresql://postgres:postgres@localhost:5432/gomigrator?sslmode=disable # DSN connection string to DB
  dir: ./migrations # folder for migration files
  type: sql # "go" or "sql"
  table_name: migrations # name of table with migrations

logger:
  level: INFO
```

Для запуска мигратора с нужным файлом конфигурации используем:

```bash
gomigrator -config="configs/config.yml"
```

По умолчанию используется путь `configs/config.yml`.

Проверяем, что все работает корректно:

```bash
gomigrator help

Usage: gomigrator [OPTIONS] COMMAND [arg...]

  You can override varuables from config file by ENV, just use something like "${DB_DSN}"

  OPTIONS:
    -config         Path to configuration file

  COMMAND:
    create [name]   Create migration with 'name'
    up              Migrate the DB to the most recent version available
    down            Roll back the version by 1
    redo            Re-run the latest migration
    status          Print all migrations status
    dbversion       Print migrations status (last applied migration)
    help            Print usage
    version         Application version

  Examples:
    gomigrator -config="../configs/config-test.yml" create "create_user_table"
    DB_DSN="postgresql://app:test@pgsql:5432/app" gomigrator up

Feel free to put PR here: https://github.com/XanderKon/sql-migrator-otus

Inspired by:
https://github.com/pressly/goose
https://github.com/golang-migrate/migrate
```

**Создание миграции**

```bash
gomigrator create test_migration

2024-01-25 00:17:07 [INFO] Success create new migration 1706131027592_test_migration.sql
```

Миграция будет создана в директории, указанной в файле конфигурации.

**Запуск всех миграции**

```bash
gomigrator up

2024-01-25 00:17:19 [INFO] Migration 1706130758469 successfully applied!
2024-01-25 00:17:19 [INFO] Migration 1706130758470 successfully applied!
2024-01-25 00:17:19 [INFO] Migration 1706130758471 successfully applied!
2024-01-25 00:17:19 [INFO] Migration 1706131027592 successfully applied!
```

**Откат последней выполненной миграции**

```bash
gomigrator down

2024-01-25 00:17:29 [INFO] Migration 1706131027592 successfully rollback!
```

**Повтор последней миграции**

```bash
gomigrator redo

2024-01-25 00:17:40 [INFO] Migration 1706130758471 successfully rollback!
2024-01-25 00:17:40 [INFO] Migration 1706130758471 successfully applied!
```

**Вывод версии базы**

```bash
gomigrator status

+---+---------------+----------------------------------------+---------------------+
| # |       VERSION | NAME                                   | APPLIED AT          |
+---+---------------+----------------------------------------+---------------------+
| 1 | 1706130758469 | 1706130758469_test_migration_go.sql    | 2024-01-25 00:17:19 |
| 2 | 1706130758470 | 1706130758470_second_migration_sql.sql | 2024-01-25 00:17:19 |
| 3 | 1706130758471 | 1706130758471_third_migration.sql      | 2024-01-25 00:17:40 |
+---+---------------+----------------------------------------+---------------------+
|   |         TOTAL | 3                                      |                     |
+---+---------------+----------------------------------------+---------------------+
```

**Вывод статуса миграций**

```bash
gomigrator status

2024-01-25 00:18:29 [INFO] Current migration version: 1706130758471
```

## Демо-режим

Для демонстрации работы приложения можно использовать команду из Make-файла:

```bash
make run-compose-demo

...

+---+---------+------+------------+
| # | VERSION | NAME | APPLIED AT |
+---+---------+------+------------+
+---+---------+------+------------+
|   |   TOTAL |    0 |            |
+---+---------+------+------------+
2024-01-24 21:24:19 [INFO] Migration 1706128932160 successfully applied!
2024-01-24 21:24:19 [INFO] Migration 1706128932162 successfully applied!
2024-01-24 21:24:19 [INFO] Current migration version: 1706128932162
+---+---------------+----------------------------------------+---------------------+
| # |       VERSION | NAME                                   | APPLIED AT          |
+---+---------------+----------------------------------------+---------------------+
| 1 | 1706128932160 | 1706128932160_test_migration_go.sql    | 2024-01-24 21:24:19 |
| 2 | 1706128932162 | 1706128932162_second_migration_sql.sql | 2024-01-24 21:24:19 |
+---+---------------+----------------------------------------+---------------------+
|   |         TOTAL | 2                                      |                     |
+---+---------------+----------------------------------------+---------------------+
2024-01-24 21:24:19 [INFO] Migration 1706128932162 successfully rollback!
+---+---------------+-------------------------------------+---------------------+
| # |       VERSION | NAME                                | APPLIED AT          |
+---+---------------+-------------------------------------+---------------------+
| 1 | 1706128932160 | 1706128932160_test_migration_go.sql | 2024-01-24 21:24:19 |
+---+---------------+-------------------------------------+---------------------+
|   |         TOTAL | 1                                   |                     |
+---+---------------+-------------------------------------+---------------------+
2024-01-24 21:24:19 [INFO] Migration 1706128932160 successfully rollback!
2024-01-24 21:24:19 [INFO] Migration 1706128932160 successfully applied!
+---+---------------+-------------------------------------+---------------------+
| # |       VERSION | NAME                                | APPLIED AT          |
+---+---------------+-------------------------------------+---------------------+
| 1 | 1706128932160 | 1706128932160_test_migration_go.sql | 2024-01-24 21:24:19 |
+---+---------------+-------------------------------------+---------------------+
|   |         TOTAL | 1                                   |                     |
+---+---------------+-------------------------------------+---------------------+
2024-01-24 21:24:19 [INFO] Success create new migration 1706131459910_test_new_migration.sql
2024-01-24 21:24:20 [INFO] Migration 1706128932162 successfully applied!
2024-01-24 21:24:20 [INFO] Migration 1706131459910 successfully applied!
+---+---------------+----------------------------------------+---------------------+
| # |       VERSION | NAME                                   | APPLIED AT          |
+---+---------------+----------------------------------------+---------------------+
| 1 | 1706128932160 | 1706128932160_test_migration_go.sql    | 2024-01-24 21:24:19 |
| 2 | 1706128932162 | 1706128932162_second_migration_sql.sql | 2024-01-24 21:24:19 |
| 3 | 1706131459910 | 1706131459910_test_new_migration.sql   | 2024-01-24 21:24:20 |
+---+---------------+----------------------------------------+---------------------+
|   |         TOTAL | 3                                      |                     |
+---+---------------+----------------------------------------+---------------------+

...
```

Команда поднимент пару контейнеров (с приложением и БД), выполнит основные команды и завершится. При этом будут использованы тестовые миграции в директории `build/migrations`, а также конфиг по умолчанию из директории `configs/config.yml`.

Задача: [MISSION.md](docs/MISSION.md)
