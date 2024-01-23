# Инструмент "SQL-мигратор"

[![Test](https://github.com/XanderKon/sql-migrator-otus/actions/workflows/actions.yml/badge.svg)](https://github.com/XanderKon/sql-migrator-otus/actions/workflows/actions.yml)

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

| Параметр | Описание                       | Возможные значения                                                       |
| -------- | ------------------------------ | ------------------------------------------------------------------------ |
| dsn      | Строка подключения к БД        | postgresql://app:!ChangeMe!@pgsql:5432/app?serverVersion=15&charset=utf8 |
| dir      | Директория для хранения файлов | migrations                                                               |
| type     | Тип генерируемых файлов        | go \| sql                                                                |

В конфигурации можно использовать переменные окружения, тогда в качестве значения нужно использовать специальную нотацию: `${ENV_VAR}` или `$ENV_VAR`.

Пример:

`configs/config.yml`

```
dsn: "$DB_DSN"
dir: "migrations"
type: "go"
```

`DB_DSN=DSN_FROM_ENV ./bin/gomigrator`

Задача: [MISSION.md](docs/MISSION.md)
